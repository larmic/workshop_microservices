package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Hausaufgaben-Check auf der Startseite: Das Dashboard startet das Image der
// Hausaufgabe kurz selbst im Compose-Netz und ruft dort GET /health auf.
// Der eigene Container des Teilnehmers muss dafür nicht laufen, Host-Port
// 8080 (Traefik) spielt keine Rolle.
//
// Der Docker-Socket ist beschreibbar gemountet. Deshalb gilt: Aus dem Request
// kommt nur eine Image-ID, die gegen die frisch erzeugte Kandidatenliste
// geprüft wird. Freier Text landet nie in einem docker-Befehl.

const (
	homeworkContainer   = "workshop-homework-check"
	homeworkLabel       = "workshop.homework-check=true"
	homeworkPreferred   = "my-booking-service"
	homeworkBudget      = 45 * time.Second
	homeworkBadStatusIn = 5 * time.Second
	workshopVendor      = "Team Neusta"
)

var imageIDPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// HomeworkImage ist ein Kandidat für die Auswahl auf der Startseite.
type HomeworkImage struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Created      string `json:"created"`
	BuiltLocally bool   `json:"builtLocally"`
	WorkshopCode bool   `json:"workshopCode"`
}

type homeworkImageList struct {
	DockerAvailable bool            `json:"dockerAvailable"`
	Preferred       string          `json:"preferred,omitempty"`
	Images          []HomeworkImage `json:"images"`
}

type homeworkCheckRequest struct {
	ImageID string `json:"imageId"`
}

// HomeworkResult: result ist einer von ok, exited, timeout, bad-status,
// start-failed, unreachable, unknown-image, choose, docker-error, busy.
// source sagt, was geprüft wurde: running (laufender Container), started
// (gestoppter, kurz gestartet) oder test (Testcontainer aus dem Image).
type HomeworkResult struct {
	Result    string `json:"result"`
	Source    string `json:"source,omitempty"`
	Container string `json:"container,omitempty"`
	Image     string `json:"image,omitempty"`
	Status    int    `json:"status,omitempty"`
	ExitCode  int    `json:"exitCode,omitempty"`
	Seconds   int    `json:"seconds"`
	Logs      string `json:"logs,omitempty"`
}

// dockerInspectImage ist der Ausschnitt aus `docker image inspect`, den der
// Filter braucht.
type dockerInspectImage struct {
	ID          string   `json:"Id"`
	RepoTags    []string `json:"RepoTags"`
	RepoDigests []string `json:"RepoDigests"`
	Created     string   `json:"Created"`
	Metadata    struct {
		LastTagTime string `json:"LastTagTime"`
	} `json:"Metadata"`
	Config struct {
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
}

// Homework bündelt die Handler; der Mutex sorgt dafür, dass immer nur ein
// Check läuft.
type Homework struct {
	composeFiles []string
	fallbackNet  string
	mu           sync.Mutex
	probe        *http.Client
}

func NewHomework(composeFiles []string, projectName string) *Homework {
	return &Homework{
		composeFiles: composeFiles,
		fallbackNet:  projectName + "_default",
		// Ohne Proxy: Firmen-Proxy-Variablen dürfen Verkehr im Compose-Netz
		// nicht umleiten.
		probe: &http.Client{Timeout: time.Second, Transport: &http.Transport{Proxy: nil}},
	}
}

func docker(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil && stderr.Len() > 0 {
		log.Printf("docker %s: %s", args[0], strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), err
}

// CleanupHomeworkContainers räumt Reste eines abgebrochenen Checks weg.
func CleanupHomeworkContainers() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := docker(ctx, "ps", "-aq", "--filter", "label="+homeworkLabel)
	if err != nil {
		return
	}
	ids := strings.Fields(out)
	if len(ids) == 0 {
		return
	}
	_, _ = docker(ctx, append([]string{"rm", "-f"}, ids...)...)
}

// ImagesHandler liefert die Kandidaten für die Image-Auswahl.
func (h *Homework) ImagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		list := homeworkImageList{Images: []HomeworkImage{}}
		images, err := h.candidates(ctx)
		if err == nil {
			list.DockerAvailable = true
			list.Images = images
			for _, img := range images {
				if repoOf(img.Name) == homeworkPreferred {
					list.Preferred = img.ID
					break
				}
			}
		}
		writeJSON(w, list)
	}
}

// CheckHandler prüft die Hausaufgabe in dieser Reihenfolge:
//  1. Läuft ein Container namens my-booking-service (oder aus dem Image)?
//     Dann diesen prüfen.
//  2. Gibt es einen gestoppten? Dann ihn starten, prüfen, wieder stoppen.
//  3. Gibt es das Image? Dann einen Testcontainer daraus prüfen.
//  4. Sonst "choose": Die Seite zeigt die Image-Auswahl.
//
// Mit imageId (nach der Auswahl) gilt dasselbe für das gewählte Image:
// erst ein laufender Container daraus, sonst ein Testcontainer.
func (h *Homework) CheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req homeworkCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "ungültiger Body", http.StatusBadRequest)
			return
		}
		if req.ImageID != "" && !imageIDPattern.MatchString(req.ImageID) {
			http.Error(w, "imageId ist ungültig", http.StatusBadRequest)
			return
		}
		if !h.mu.TryLock() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(HomeworkResult{Result: "busy"})
			return
		}
		defer h.mu.Unlock()

		listCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		images, errImages := h.candidates(listCtx)
		containers, errContainers := h.containers(listCtx)
		cancel()
		if errImages != nil || errContainers != nil {
			writeJSON(w, HomeworkResult{Result: "docker-error"})
			return
		}

		var image *HomeworkImage
		if req.ImageID != "" {
			for i := range images {
				if images[i].ID == req.ImageID {
					image = &images[i]
				}
			}
			if image == nil {
				writeJSON(w, HomeworkResult{Result: "unknown-image"})
				return
			}
		} else {
			for i := range images {
				if repoOf(images[i].Name) == homeworkPreferred {
					image = &images[i]
				}
			}
		}

		imageID := ""
		if image != nil {
			imageID = image.ID
		}
		running, stopped := matchHomeworkContainers(containers, imageID, req.ImageID == "")

		var res HomeworkResult
		switch {
		case running != nil:
			log.Printf("Homework check: running container %s", running.name())
			res = h.checkExisting(r.Context(), *running, false)
		case stopped != nil && req.ImageID == "":
			log.Printf("Homework check: stopped container %s", stopped.name())
			res = h.checkExisting(r.Context(), *stopped, true)
		case image != nil:
			log.Printf("Homework check: test container from %s", image.Name)
			res = h.run(r.Context(), image.ID)
			res.Image = image.Name
		default:
			res = HomeworkResult{Result: "choose"}
		}
		writeJSON(w, res)
	}
}

// dockerContainer ist der Ausschnitt aus `docker inspect` für Container.
type dockerContainer struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Image  string `json:"Image"`
	Config struct {
		Image  string            `json:"Image"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
	State struct {
		Running bool `json:"Running"`
	} `json:"State"`
	NetworkSettings struct {
		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

func (c dockerContainer) name() string { return strings.TrimPrefix(c.Name, "/") }

func (h *Homework) containers(ctx context.Context) ([]dockerContainer, error) {
	out, err := docker(ctx, "ps", "-aq", "--no-trunc")
	if err != nil {
		return nil, err
	}
	ids := uniqueFields(out)
	if len(ids) == 0 {
		return nil, nil
	}
	out, err = docker(ctx, append([]string{"inspect"}, ids...)...)
	if err != nil {
		return nil, err
	}
	var containers []dockerContainer
	if err := json.Unmarshal([]byte(out), &containers); err != nil {
		return nil, err
	}
	return containers, nil
}

// matchHomeworkContainers sucht den Container der Hausaufgabe: aus dem Image
// imageID oder, wenn byName gesetzt ist, mit Namen bzw. Image
// my-booking-service. Container des Workshop-Stacks (Compose-Label) und der
// eigene Testcontainer zählen nie.
func matchHomeworkContainers(containers []dockerContainer, imageID string, byName bool) (running, stopped *dockerContainer) {
	for i := range containers {
		c := &containers[i]
		if c.Config.Labels["com.docker.compose.project"] != "" || c.name() == homeworkContainer {
			continue
		}
		match := imageID != "" && c.Image == imageID
		if byName && (c.name() == homeworkPreferred || repoOf(c.Config.Image) == homeworkPreferred) {
			match = true
		}
		if !match {
			continue
		}
		if c.State.Running && running == nil {
			running = c
		}
		if !c.State.Running && stopped == nil {
			stopped = c
		}
	}
	return running, stopped
}

// checkExisting prüft einen vorhandenen Container. Er hängt meist im Netz
// bridge und ist vom Dashboard aus nicht erreichbar, deshalb kommt er kurz
// zusätzlich ins Compose-Netz. Ein gestoppter wird gestartet und danach
// wieder gestoppt. Klappt das nicht (z. B. Host-Port 8080 belegt, weil
// Traefik ihn hat), gibt es einen Testcontainer aus seinem Image.
func (h *Homework) checkExisting(ctx context.Context, c dockerContainer, wasStopped bool) HomeworkResult {
	start := time.Now()
	bg := func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), 15*time.Second)
	}
	fallback := func() HomeworkResult {
		res := h.run(ctx, c.Image)
		res.Image = c.Config.Image
		return res
	}

	if wasStopped {
		sctx, cancel := bg()
		_, err := docker(sctx, "start", c.ID)
		cancel()
		if err != nil {
			return fallback()
		}
		defer func() {
			sctx, cancel := bg()
			defer cancel()
			_, _ = docker(sctx, "stop", c.ID)
		}()
	}

	netCtx, cancel := bg()
	network := h.network(netCtx)
	cancel()
	if _, ok := c.NetworkSettings.Networks[network]; !ok {
		cctx, cancel := bg()
		_, err := docker(cctx, "network", "connect", network, c.ID)
		cancel()
		if err != nil {
			if wasStopped {
				return fallback()
			}
			return HomeworkResult{Result: "unreachable", Container: c.name(), Seconds: int(time.Since(start).Seconds())}
		}
		defer func() {
			dctx, cancel := bg()
			defer cancel()
			_, _ = docker(dctx, "network", "disconnect", network, c.ID)
		}()
	}

	ictx, cancel := bg()
	ip, err := docker(ictx, "inspect", "-f", "{{(index .NetworkSettings.Networks \""+network+"\").IPAddress}}", c.ID)
	cancel()
	ip = strings.TrimSpace(ip)
	if err != nil || net.ParseIP(ip) == nil {
		return HomeworkResult{Result: "unreachable", Container: c.name(), Seconds: int(time.Since(start).Seconds())}
	}

	budget := 15 * time.Second
	source := "running"
	if wasStopped {
		budget = homeworkBudget
		source = "started"
	}
	res := h.waitHealthy(ctx, c.ID, net.JoinHostPort(ip, "8080"), start, budget)
	res.Source = source
	res.Container = c.name()
	res.Image = c.Config.Image
	return res
}

func (h *Homework) run(ctx context.Context, imageID string) HomeworkResult {
	start := time.Now()

	// Aufräumen läuft mit eigenem Context, damit ein abgebrochener Request
	// (Tab zu, Reload) keinen Container stehen lässt.
	cleanup := func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = docker(c, "rm", "-f", homeworkContainer)
	}
	cleanup()
	defer cleanup()

	runCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	// Ohne --rm: Nach einem Absturz sollen Exit-Code und Logs lesbar bleiben.
	_, err := docker(runCtx, "run", "-d", "--pull", "never",
		"--name", homeworkContainer, "--label", homeworkLabel,
		"--network", h.network(runCtx),
		"--memory", "768m", "--cpus", "1", "--pids-limit", "512",
		imageID)
	cancel()
	if err != nil {
		return HomeworkResult{Result: "start-failed", Source: "test", Seconds: int(time.Since(start).Seconds()), Logs: h.logs(homeworkContainer)}
	}
	res := h.waitHealthy(ctx, homeworkContainer, homeworkContainer+":8080", start, homeworkBudget)
	res.Source = "test"
	return res
}

// waitHealthy ruft GET /health auf, bis 200 kommt, der Container sich
// beendet, ein anderer Status 5 s lang bleibt oder das Budget abläuft.
func (h *Homework) waitHealthy(ctx context.Context, container, hostPort string, start time.Time, budget time.Duration) HomeworkResult {
	seconds := func() int { return int(time.Since(start).Seconds()) }
	deadline := start.Add(budget)
	lastStatus := 0
	var badSince time.Time
	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			return HomeworkResult{Result: "timeout", Seconds: seconds()}
		}
		if exited, code := h.exited(container); exited {
			return HomeworkResult{Result: "exited", ExitCode: code, Seconds: seconds(), Logs: h.logs(container)}
		}
		resp, err := h.probe.Get("http://" + hostPort + "/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return HomeworkResult{Result: "ok", Status: http.StatusOK, Seconds: seconds()}
			}
			if lastStatus == 0 {
				badSince = time.Now()
			}
			lastStatus = resp.StatusCode
			if time.Since(badSince) >= homeworkBadStatusIn {
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	if lastStatus != 0 {
		return HomeworkResult{Result: "bad-status", Status: lastStatus, Seconds: seconds()}
	}
	return HomeworkResult{Result: "timeout", Seconds: seconds(), Logs: h.logs(container)}
}

func (h *Homework) exited(container string) (bool, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := docker(ctx, "inspect", "-f", "{{.State.Status}} {{.State.ExitCode}}", container)
	if err != nil {
		return false, 0
	}
	var status string
	var code int
	if _, err := fmt.Sscan(out, &status, &code); err != nil {
		return false, 0
	}
	return status == "exited" || status == "dead", code
}

func (h *Homework) logs(container string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "logs", "--tail", "20", container)
	out, _ := cmd.CombinedOutput()
	return strings.TrimSpace(string(out))
}

// network leitet das Compose-Netz vom eigenen Container ab (Hostname =
// Container-ID). Fallback ist <Projekt>_default.
func (h *Homework) network(ctx context.Context) string {
	host, err := os.Hostname()
	if err != nil {
		return h.fallbackNet
	}
	out, err := docker(ctx, "inspect", "-f", "{{range $k, $v := .NetworkSettings.Networks}}{{$k}} {{end}}", host)
	if err != nil {
		return h.fallbackNet
	}
	if nets := strings.Fields(out); len(nets) > 0 {
		return nets[0]
	}
	return h.fallbackNet
}

func (h *Homework) candidates(ctx context.Context) ([]HomeworkImage, error) {
	out, err := docker(ctx, "image", "ls", "-q", "--no-trunc")
	if err != nil {
		return nil, err
	}
	ids := uniqueFields(out)
	if len(ids) == 0 {
		return []HomeworkImage{}, nil
	}
	out, err = docker(ctx, append([]string{"image", "inspect"}, ids...)...)
	if err != nil {
		return nil, err
	}
	var inspected []dockerInspectImage
	if err := json.Unmarshal([]byte(out), &inspected); err != nil {
		return nil, err
	}
	return filterHomeworkImages(inspected, composeImageRepos(h.composeFiles)), nil
}

// filterHomeworkImages entfernt alles, was zum Workshop gehört: Images
// namens workshop-microservices-* (Docker Hub und lokal gebaut) sowie die
// Infrastruktur-Images aus den Compose-Dateien. Images mit dem Vendor-Label
// des Workshops, aber anderem Namen (kopierte Beispiel-Lösung) bleiben
// sichtbar und werden markiert. Selbst gebaute zuerst, dann die neuesten.
func filterHomeworkImages(images []dockerInspectImage, infraRepos map[string]bool) []HomeworkImage {
	result := []HomeworkImage{}
	for _, img := range images {
		name := ""
		skip := false
		for _, tag := range img.RepoTags {
			if tag == "" || tag == "<none>:<none>" || strings.Contains(tag, "@") {
				continue
			}
			repo := repoOf(tag)
			if strings.Contains(repo, "workshop-microservices-") || infraRepos[repo] {
				skip = true
			}
			if name == "" || repo == homeworkPreferred {
				name = tag
			}
		}
		if skip || name == "" {
			continue
		}
		result = append(result, HomeworkImage{
			ID:           img.ID,
			Name:         name,
			Created:      img.Created,
			BuiltLocally: builtLocally(img),
			WorkshopCode: img.Config.Labels["org.opencontainers.image.vendor"] == workshopVendor,
		})
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].BuiltLocally != result[j].BuiltLocally {
			return result[i].BuiltLocally
		}
		return result[i].Created > result[j].Created
	})
	return result
}

// builtLocally erkennt selbst gebaute Images. Im klassischen Image-Store
// haben sie keinen Registry-Digest. Im containerd-Store (Standard bei neueren
// Docker-Desktop-Installationen) haben auch lokale Builds einen; dort hilft
// die Zeit: Ein lokaler Build wird direkt nach dem Erzeugen getaggt, ein
// gepulltes Image trägt das Build-Datum des Herausgebers.
func builtLocally(img dockerInspectImage) bool {
	if len(img.RepoDigests) == 0 {
		return true
	}
	created, err1 := time.Parse(time.RFC3339Nano, img.Created)
	tagged, err2 := time.Parse(time.RFC3339Nano, img.Metadata.LastTagTime)
	if err1 != nil || err2 != nil {
		return false
	}
	gap := tagged.Sub(created)
	return gap >= -time.Minute && gap < 15*time.Minute
}

var composeImageLine = regexp.MustCompile(`^\s*image:\s*["']?([^"'\s]+)`)

// composeImageRepos liest die image:-Zeilen der Compose-Dateien. So fallen
// neue Infrastruktur-Images automatisch aus der Auswahl.
func composeImageRepos(files []string) map[string]bool {
	repos := map[string]bool{}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if m := composeImageLine.FindStringSubmatch(scanner.Text()); m != nil {
				repos[repoOf(m[1])] = true
			}
		}
		_ = file.Close()
	}
	return repos
}

// repoOf schneidet Tag und Digest ab: "traefik:v3.7" -> "traefik".
func repoOf(ref string) string {
	if i := strings.Index(ref, "@"); i >= 0 {
		ref = ref[:i]
	}
	if i := strings.LastIndex(ref, ":"); i > strings.LastIndex(ref, "/") {
		ref = ref[:i]
	}
	return ref
}

func uniqueFields(s string) []string {
	seen := map[string]bool{}
	var out []string
	for _, f := range strings.Fields(s) {
		if !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	return out
}

// LocalOriginOnly lässt nur Aufrufe der eigenen Seite durch. Das Dashboard
// setzt sonst Access-Control-Allow-Origin: *, damit könnte jede offene
// Webseite die Image-Namen lesen oder einen Check auslösen. Erlaubt ist ein
// fehlender Origin oder einer, dessen Host gleich dem Host-Header ist und auf
// eine lokale oder private Adresse zeigt (schützt auch vor DNS-Rebinding).
func LocalOriginOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !originAllowed(r.Header.Get("Origin"), r.Host) {
			http.Error(w, "nur von localhost erlaubt", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func originAllowed(origin, host string) bool {
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	if !strings.EqualFold(normalizeHost(u.Host, u.Scheme), normalizeHost(host, u.Scheme)) {
		return false
	}
	return isLocalHost(u.Hostname())
}

func normalizeHost(host, scheme string) string {
	if scheme == "http" {
		host = strings.TrimSuffix(host, ":80")
	}
	if scheme == "https" {
		host = strings.TrimSuffix(host, ":443")
	}
	return host
}

func isLocalHost(name string) bool {
	name = strings.ToLower(name)
	if name == "localhost" || strings.HasSuffix(name, ".localhost") {
		return true
	}
	ip := net.ParseIP(name)
	return ip != nil && (ip.IsLoopback() || ip.IsPrivate())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
