package consul

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type serviceEntry struct {
	Service struct {
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

type Resolver struct {
	consulURL string
	client    *http.Client
}

func NewResolver(consulURL string, client *http.Client) *Resolver {
	return &Resolver{
		consulURL: consulURL,
		client:    client,
	}
}

func (r *Resolver) ResolveServiceURL(serviceName string) (string, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s?passing=true", r.consulURL, serviceName)

	resp, err := r.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("consul request failed: %w", err)
	}
	defer resp.Body.Close()

	var entries []serviceEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return "", fmt.Errorf("consul response parse failed: %w", err)
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no healthy instances found for %s", serviceName)
	}

	entry := entries[rand.Intn(len(entries))]
	return fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port), nil
}
