package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

func Register(consulURL string, cfg ServiceConfig) (string, error) {
	const maxRetries = 5
	initialBackoff := 1 * time.Second

	hostname, err := os.Hostname()
	if err != nil {
		hostname = cfg.Address
	}

	serviceID := fmt.Sprintf("%s-%s", cfg.Name, hostname)

	payload := registrationPayload{
		Name:    cfg.Name,
		ID:      serviceID,
		Address: hostname,
		Port:    cfg.Port,
		Check: healthCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, cfg.Port),
			Interval: "10s",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registration payload: %w", err)
	}

	var lastErr error
	backoff := initialBackoff

	for attempt := 1; attempt <= maxRetries; attempt++ {
		lastErr = registerOnce(consulURL, body)
		if lastErr == nil {
			log.Printf("Registered service %q (id=%s) at %s:%d with Consul", cfg.Name, serviceID, hostname, cfg.Port)
			return serviceID, nil
		}

		if attempt < maxRetries {
			log.Printf("Consul registration attempt %d/%d failed: %v. Retrying in %v...", attempt, maxRetries, lastErr, backoff)
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return "", fmt.Errorf("consul registration failed after %d attempts: %w", maxRetries, lastErr)
}

func registerOnce(consulURL string, body []byte) error {
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
