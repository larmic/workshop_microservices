package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type ChaosState struct {
	Mode      string `json:"mode"`
	LatencyMs int    `json:"latencyMs"`
}

type InstanceInfo struct {
	consul.Instance
	Chaos     *ChaosState `json:"chaos,omitempty"`
	Reachable bool        `json:"reachable"`
}

type SetChaosRequest struct {
	Mode        string   `json:"mode"`
	InstanceIDs []string `json:"instanceIds,omitempty"`
}

type SetChaosResult struct {
	InstanceID string      `json:"instanceId"`
	OK         bool        `json:"ok"`
	Error      string      `json:"error,omitempty"`
	Chaos      *ChaosState `json:"chaos,omitempty"`
}

func consulName(short string) string { return short + "-service" }

func ListInstancesHandler(resolver *consul.Resolver, allowedServices []string) http.HandlerFunc {
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

		targetSet := make(map[string]bool, len(req.InstanceIDs))
		for _, id := range req.InstanceIDs {
			targetSet[id] = true
		}

		body, _ := json.Marshal(struct {
			Mode string `json:"mode"`
		}{Mode: req.Mode})

		client := &http.Client{Timeout: 2 * time.Second}
		results := make([]SetChaosResult, 0, len(instances))
		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, inst := range instances {
			if len(targetSet) > 0 && !targetSet[inst.ID] {
				continue
			}
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
