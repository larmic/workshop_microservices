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

func HealthOverviewHandler(composeArgs []string, scalableServices []string, resolver *consul.Resolver, bookingURLs map[string]string) http.HandlerFunc {
	bookingNames := make([]string, 0, len(bookingURLs))
	for name := range bookingURLs {
		bookingNames = append(bookingNames, name)
	}
	sort.Strings(bookingNames)

	return func(w http.ResponseWriter, r *http.Request) {
		results := make([]ServiceHealth, len(bookingNames)+len(scalableServices))
		client := &http.Client{Timeout: 800 * time.Millisecond}

		var wg sync.WaitGroup
		for i, name := range bookingNames {
			wg.Add(1)
			go func(i int, name string) {
				defer wg.Done()
				results[i] = ServiceHealth{
					Name:   name,
					Kind:   "booking",
					Status: probeBookingHealth(client, bookingURLs[name]),
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
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthOverview{Services: results})
	}
}

func probeBookingHealth(client *http.Client, baseURL string) string {
	resp, err := client.Get(baseURL + "/health")
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
