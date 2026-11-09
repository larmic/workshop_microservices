package handler

import (
	"os"
	"path/filepath"
	"testing"
)

func img(id string, tags []string, digests []string, created, vendor string) dockerInspectImage {
	i := dockerInspectImage{ID: id, RepoTags: tags, RepoDigests: digests, Created: created}
	if vendor != "" {
		i.Config.Labels = map[string]string{"org.opencontainers.image.vendor": vendor}
	}
	return i
}

func TestFilterHomeworkImages(t *testing.T) {
	infra := map[string]bool{"traefik": true, "hashicorp/consul": true}
	images := []dockerInspectImage{
		img("hub", []string{"larmic/workshop-microservices-booking:story3"}, []string{"d"}, "2026-11-01", workshopVendor),
		img("local-ref", []string{"workshop-microservices-booking:story1"}, nil, "2026-11-02", workshopVendor),
		img("sample", []string{"larmic/workshop-microservices-booking:custom"}, nil, "2026-11-03", workshopVendor),
		img("traefik", []string{"traefik:v3.7"}, []string{"d"}, "2026-10-01", ""),
		img("dangling", []string{"<none>:<none>"}, nil, "2026-11-04", ""),
		img("postgres", []string{"postgres:16"}, []string{"d"}, "2026-11-05", ""),
		img("old", []string{"spring-hello:latest"}, nil, "2026-09-01", ""),
		img("new", []string{"booking-go:latest"}, nil, "2026-11-06", ""),
		img("copied", []string{"my-copy:latest"}, nil, "2026-10-15", workshopVendor),
	}

	got := filterHomeworkImages(images, infra)

	want := []string{"new", "copied", "old", "postgres"}
	if len(got) != len(want) {
		t.Fatalf("got %d images %+v, want %v", len(got), got, want)
	}
	for i, id := range want {
		if got[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, got[i].ID, id)
		}
	}
	if !got[1].WorkshopCode {
		t.Error("copied sample should be marked as workshop code")
	}
	if got[3].BuiltLocally {
		t.Error("pulled image should not count as built locally")
	}
}

func TestFilterPrefersMyBookingServiceTag(t *testing.T) {
	got := filterHomeworkImages([]dockerInspectImage{
		img("x", []string{"other:1", "my-booking-service:latest"}, nil, "2026-11-01", ""),
	}, nil)
	if len(got) != 1 || got[0].Name != "my-booking-service:latest" {
		t.Fatalf("got %+v", got)
	}
}

func TestBuiltLocally(t *testing.T) {
	local := img("a", []string{"booking-go:latest"}, []string{"booking-go@sha256:1"}, "2026-11-09T12:11:30.2Z", "")
	local.Metadata.LastTagTime = "2026-11-09T12:11:30.3Z"
	pulled := img("b", []string{"python:3.13"}, []string{"python@sha256:2"}, "2026-11-05T19:30:11Z", "")
	pulled.Metadata.LastTagTime = "2026-11-09T12:11:30Z"
	classic := img("c", []string{"x:1"}, nil, "2026-01-01T00:00:00Z", "")

	if !builtLocally(local) {
		t.Error("containerd store: image tagged right after creation should count as built locally")
	}
	if builtLocally(pulled) {
		t.Error("image created days before tagging should count as pulled")
	}
	if !builtLocally(classic) {
		t.Error("classic store: image without digest should count as built locally")
	}
}

func TestRepoOf(t *testing.T) {
	cases := map[string]string{
		"traefik:v3.7":                "traefik",
		"hashicorp/consul:2.1":        "hashicorp/consul",
		"localhost:5000/team/app:1.0": "localhost:5000/team/app",
		"localhost:5000/team/app":     "localhost:5000/team/app",
		"nginx@sha256:abc":            "nginx",
		"my-booking-service":          "my-booking-service",
	}
	for in, want := range cases {
		if got := repoOf(in); got != want {
			t.Errorf("repoOf(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestComposeImageRepos(t *testing.T) {
	f := filepath.Join(t.TempDir(), "compose.yml")
	content := "services:\n  traefik:\n    image: traefik:v3.7\n  ui:\n    image: \"swaggerapi/swagger-ui:v5\"\n"
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	repos := composeImageRepos([]string{f, "/does/not/exist"})
	if !repos["traefik"] || !repos["swaggerapi/swagger-ui"] || len(repos) != 2 {
		t.Fatalf("got %v", repos)
	}
}

func TestOriginAllowed(t *testing.T) {
	cases := []struct {
		origin, host string
		want         bool
	}{
		{"", "localhost", true},
		{"http://localhost", "localhost", true},
		{"http://localhost:80", "localhost", true},
		{"http://127.0.0.1", "127.0.0.1", true},
		{"http://[::1]", "[::1]", true},
		{"http://172.28.1.5", "172.28.1.5", true},
		{"http://evil.example", "localhost", false},
		{"http://evil.example", "evil.example", false},
		{"http://localhost:3000", "localhost", false},
		{"null", "localhost", false},
	}
	for _, c := range cases {
		if got := originAllowed(c.origin, c.host); got != c.want {
			t.Errorf("originAllowed(%q, %q) = %v, want %v", c.origin, c.host, got, c.want)
		}
	}
}
