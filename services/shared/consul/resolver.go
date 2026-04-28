package consul

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type serviceEntry struct {
	Service struct {
		ID      string `json:"ID"`
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

type Instance struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	URL     string `json:"url"`
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
	instances, err := r.ResolveAllServiceURLs(serviceName)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("no healthy instances found for %s", serviceName)
	}
	return instances[rand.Intn(len(instances))].URL, nil
}

func (r *Resolver) ResolveAllServiceURLs(serviceName string) ([]Instance, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s?passing=true", r.consulURL, serviceName)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("consul request failed: %w", err)
	}
	defer resp.Body.Close()

	var entries []serviceEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("consul response parse failed: %w", err)
	}

	instances := make([]Instance, 0, len(entries))
	for _, e := range entries {
		instances = append(instances, Instance{
			ID:      e.Service.ID,
			Address: e.Service.Address,
			Port:    e.Service.Port,
			URL:     fmt.Sprintf("http://%s:%d", e.Service.Address, e.Service.Port),
		})
	}
	return instances, nil
}
