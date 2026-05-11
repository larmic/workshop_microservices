package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type ServiceHealth struct {
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Status          string `json:"status"`
	Replicas        int    `json:"replicas,omitempty"`
	HealthyReplicas int    `json:"healthyReplicas,omitempty"`
}

type HealthOverview struct {
	Services []ServiceHealth `json:"services"`
}

type InfraTarget struct {
	Name string
	URL  string
}

// HealthcheckUserAgent: damit lassen sich Probes in nginx-access-Logs
// (z. B. swagger-ui) gezielt herausfiltern.
const HealthcheckUserAgent = "workshop-dashboard-healthcheck"

// Infra-Status wird gecacht, damit das schnelle Frontend-Polling die
// Infrastruktur-Container nicht zumüllt. Booking-/scalable-Services
// bleiben uncached, weil dort Statuswechsel sofort sichtbar sein müssen.
const infraCacheTTL = 30 * time.Second

type infraCache struct {
	mu      sync.Mutex
	entries map[string]infraCacheEntry
}

type infraCacheEntry struct {
	status string
	at     time.Time
}

func (c *infraCache) get(name string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[name]
	if !ok || time.Since(e.at) > infraCacheTTL {
		return "", false
	}
	return e.status, true
}

func (c *infraCache) set(name, status string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[name] = infraCacheEntry{status: status, at: time.Now()}
}

func HealthOverviewHandler(composeArgs []string, scalableServices []string, resolver *consul.Resolver, bookingURLs map[string]string, infraTargets []InfraTarget) http.HandlerFunc {
	bookingNames := make([]string, 0, len(bookingURLs))
	for name := range bookingURLs {
		bookingNames = append(bookingNames, name)
	}
	sort.Strings(bookingNames)

	cache := &infraCache{entries: make(map[string]infraCacheEntry)}

	return func(w http.ResponseWriter, r *http.Request) {
		results := make([]ServiceHealth, len(bookingNames)+len(scalableServices)+len(infraTargets))
		client := &http.Client{Timeout: 800 * time.Millisecond}

		var wg sync.WaitGroup
		for i, name := range bookingNames {
			wg.Add(1)
			go func(i int, name string) {
				defer wg.Done()
				results[i] = ServiceHealth{
					Name:   name,
					Kind:   "booking",
					Status: probeURL(client, bookingURLs[name]+"/health"),
				}
			}(i, name)
		}

		for j, name := range scalableServices {
			idx := len(bookingNames) + j
			wg.Add(1)
			go func(idx int, name string) {
				defer wg.Done()
				results[idx] = probeScalableHealth(name, composeArgs, resolver)
			}(idx, name)
		}

		for k, t := range infraTargets {
			idx := len(bookingNames) + len(scalableServices) + k
			wg.Add(1)
			go func(idx int, t InfraTarget) {
				defer wg.Done()
				status, ok := cache.get(t.Name)
				if !ok {
					status = probeURL(client, t.URL)
					cache.set(t.Name, status)
				}
				results[idx] = ServiceHealth{Name: t.Name, Kind: "infra", Status: status}
			}(idx, t)
		}
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthOverview{Services: results})
	}
}

func probeURL(client *http.Client, url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "starting"
	}
	req.Header.Set("User-Agent", HealthcheckUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "starting"
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return "ready"
	}
	return "starting"
}

func probeScalableHealth(name string, composeArgs []string, resolver *consul.Resolver) ServiceHealth {
	sh := ServiceHealth{Name: name, Kind: "scalable"}
	sh.Replicas = countRunningContainers(composeArgs, name)
	if sh.Replicas == 0 {
		sh.Status = "scaled-to-zero"
		return sh
	}
	instances, err := resolver.ResolveAllServiceURLs(consulName(name))
	if err != nil {
		sh.Status = "starting"
		return sh
	}
	sh.HealthyReplicas = len(instances)
	if sh.HealthyReplicas >= sh.Replicas {
		sh.Status = "ready"
	} else if sh.HealthyReplicas > 0 {
		sh.Status = "starting"
	} else {
		sh.Status = "starting"
	}
	return sh
}
