package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ServiceConfig struct {
	Name    string
	Address string
	Port    int
}

type registrationPayload struct {
	Name    string      `json:"Name"`
	ID      string      `json:"ID"`
	Address string      `json:"Address"`
	Port    int         `json:"Port"`
	Check   healthCheck `json:"Check"`
}

type healthCheck struct {
	HTTP     string `json:"HTTP"`
	Interval string `json:"Interval"`
}

func Register(consulURL string, cfg ServiceConfig) error {
	payload := registrationPayload{
		Name:    cfg.Name,
		ID:      cfg.Name,
		Address: cfg.Address,
		Port:    cfg.Port,
		Check: healthCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", cfg.Address, cfg.Port),
			Interval: "10s",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal registration payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/agent/service/register", consulURL)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("consul registration request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("consul registration returned status %d", resp.StatusCode)
	}

	log.Printf("Registered service %q at %s:%d with Consul", cfg.Name, cfg.Address, cfg.Port)
	return nil
}

func Deregister(consulURL string, serviceID string) {
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", consulURL, serviceID)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Printf("WARNING: failed to create deregistration request: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("WARNING: consul deregistration request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Deregistered service %q from Consul", serviceID)
}
