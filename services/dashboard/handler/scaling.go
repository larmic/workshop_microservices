package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type ServiceInfo struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
}

type ScaleRequest struct {
	Replicas int `json:"replicas"`
}

type composeContainer struct {
	State string `json:"State"`
}

func ListServicesHandler(composeArgs []string, allowedServices []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		services := make([]ServiceInfo, 0, len(allowedServices))
		for _, name := range allowedServices {
			count := countRunningContainers(composeArgs, name)
			services = append(services, ServiceInfo{Name: name, Replicas: count})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}

func ScaleServiceHandler(composeArgs []string, allowedServices []string) http.HandlerFunc {
	allowed := make(map[string]bool, len(allowedServices))
	for _, s := range allowedServices {
		allowed[s] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		log.Printf("%s %s (service=%s) from %s", r.Method, r.URL.Path, name, r.RemoteAddr)

		if !allowed[name] {
			http.Error(w, fmt.Sprintf("Service %q is not scalable. Allowed: %v", name, allowedServices), http.StatusBadRequest)
			return
		}

		var req ScaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.Replicas < 0 || req.Replicas > 10 {
			http.Error(w, "Replicas must be between 0 and 10", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		args := append([]string{}, composeArgs...)
		args = append(args, "up", "-d", "--scale", fmt.Sprintf("%s=%d", name, req.Replicas), "--no-recreate", name)

		cmd := exec.CommandContext(ctx, "docker", append([]string{"compose"}, args...)...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		log.Printf("Executing: docker compose %s", strings.Join(args, " "))

		if err := cmd.Run(); err != nil {
			log.Printf("Scale error: %v, stderr: %s", err, stderr.String())
			http.Error(w, "Scaling failed: "+stderr.String(), http.StatusInternalServerError)
			return
		}

		count := countRunningContainers(composeArgs, name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ServiceInfo{Name: name, Replicas: count})
	}
}

func countRunningContainers(composeArgs []string, serviceName string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := append([]string{}, composeArgs...)
	args = append(args, "ps", "--format", "json", serviceName)

	cmd := exec.CommandContext(ctx, "docker", append([]string{"compose"}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to list containers for %s: %v, stderr: %s", serviceName, err, stderr.String())
		return 0
	}

	count := 0
	for _, line := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		if line == "" {
			continue
		}
		var container composeContainer
		if err := json.Unmarshal([]byte(line), &container); err != nil {
			log.Printf("Failed to parse container JSON: %v", err)
			continue
		}
		if container.State == "running" {
			count++
		}
	}
	return count
}
