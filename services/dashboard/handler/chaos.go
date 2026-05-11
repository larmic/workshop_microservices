package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type ChaosState struct {
	Mode      string     `json:"mode"`
	LatencyMs int        `json:"latencyMs"`
	// LastSeenAt: Zeitpunkt des letzten regulären (nicht-whitelisted) Requests
	// gegen diese Instanz. Vom Backend gepflegt (shared/chaos Middleware).
	// Optional, weil eine frisch gestartete Instanz noch nichts gesehen hat.
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
}

type InstanceInfo struct {
	consul.Instance
	Chaos     *ChaosState `json:"chaos,omitempty"`
	Reachable bool        `json:"reachable"`
	// ContainerName ist der docker-compose-Container-Name (z. B.
	// "services-car-4"), den der User aus den Logs kennt. Optional, weil
	// das Dashboard außerhalb einer Compose-Umgebung laufen könnte.
	ContainerName string `json:"containerName,omitempty"`
}

type SetChaosRequest struct {
	Mode string `json:"mode"`
}

type SetChaosResult struct {
	InstanceID string      `json:"instanceId"`
	OK         bool        `json:"ok"`
	Error      string      `json:"error,omitempty"`
	Chaos      *ChaosState `json:"chaos,omitempty"`
}

func consulName(short string) string { return short + "-service" }

func ListInstancesHandler(resolver *consul.Resolver, allowedServices []string, composeArgs []string) http.HandlerFunc {
	allowed := allowedSet(allowedServices)
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !allowed[name] {
			http.Error(w, fmt.Sprintf("service %q not allowed", name), http.StatusBadRequest)
			return
		}

		instances, err := resolver.ResolveAllServiceURLs(consulName(name))
		if err != nil {
			http.Error(w, fmt.Sprintf("consul lookup failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Parallel zur Backend-Befragung holen wir die docker-compose-
		// Container-Namen — der User erkennt "services-car-4" leichter
		// als die Consul-ID "car-service-ed9cda40b84d".
		composeNames := ContainerNamesByHostname(composeArgs, name)

		client := &http.Client{Timeout: 1 * time.Second}
		results := make([]InstanceInfo, len(instances))

		var wg sync.WaitGroup
		for i, inst := range instances {
			wg.Add(1)
			go func(i int, inst consul.Instance) {
				defer wg.Done()
				info := InstanceInfo{Instance: inst}
				state, err := fetchChaos(client, inst.URL)
				if err == nil {
					info.Chaos = state
					info.Reachable = true
				} else {
					info.Reachable = false
				}
				// Hostname ist das Suffix nach dem letzten "-" in der
				// Service-ID (siehe shared/consul/register.go).
				if idx := strings.LastIndex(inst.ID, "-"); idx >= 0 && idx+1 < len(inst.ID) {
					hostname := inst.ID[idx+1:]
					if cn, ok := composeNames[hostname]; ok {
						info.ContainerName = cn
					}
				}
				results[i] = info
			}(i, inst)
		}
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(results)
	}
}

func SetChaosHandler(resolver *consul.Resolver, allowedServices []string) http.HandlerFunc {
	allowed := allowedSet(allowedServices)
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !allowed[name] {
			http.Error(w, fmt.Sprintf("service %q not allowed", name), http.StatusBadRequest)
			return
		}

		var req SetChaosRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		instances, err := resolver.ResolveAllServiceURLs(consulName(name))
		if err != nil {
			http.Error(w, fmt.Sprintf("consul lookup failed: %v", err), http.StatusInternalServerError)
			return
		}

		body, _ := json.Marshal(struct {
			Mode string `json:"mode"`
		}{Mode: req.Mode})

		client := &http.Client{Timeout: 2 * time.Second}
		results := make([]SetChaosResult, 0, len(instances))
		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, inst := range instances {
			wg.Add(1)
			go func(inst consul.Instance) {
				defer wg.Done()
				res := SetChaosResult{InstanceID: inst.ID}
				state, err := postChaos(client, inst.URL, body)
				if err != nil {
					res.OK = false
					res.Error = err.Error()
				} else {
					res.OK = true
					res.Chaos = state
				}
				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}(inst)
		}
		wg.Wait()

		log.Printf("Chaos set on %s: mode=%s targets=%d", name, req.Mode, len(results))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(results)
	}
}

func fetchChaos(client *http.Client, baseURL string) (*ChaosState, error) {
	resp, err := client.Get(baseURL + "/admin/chaos")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var s ChaosState
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func postChaos(client *http.Client, baseURL string, body []byte) (*ChaosState, error) {
	resp, err := client.Post(baseURL+"/admin/chaos", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(msg))
	}
	var s ChaosState
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func allowedSet(allowed []string) map[string]bool {
	m := make(map[string]bool, len(allowed))
	for _, s := range allowed {
		m[s] = true
	}
	return m
}
