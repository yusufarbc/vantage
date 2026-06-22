// Package dockerstatus reads container health from the docker-socket-proxy
// sidecar (see the "docker-proxy" service in docker-compose.yml). It never
// talks to /var/run/docker.sock directly — only the locked-down proxy, which
// only allows GET /containers/*.
package dockerstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Container is a simplified view of a single container's health, derived
// from the Docker Engine API's container-list response.
type Container struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	State  string `json:"state"`  // e.g. "running", "exited"
	Status string `json:"status"` // e.g. "Up 21 hours (healthy)"
	Health string `json:"health"` // "healthy" | "unhealthy" | "starting" | "none"
}

type rawContainer struct {
	Names  []string `json:"Names"`
	Image  string   `json:"Image"`
	State  string   `json:"State"`
	Status string   `json:"Status"`
}

var healthPattern = regexp.MustCompile(`\(([a-z]+)\)`)

// List fetches the current state of every container visible to the
// docker-proxy sidecar at proxyURL (e.g. "http://docker-proxy:2375").
// Returns an error if proxyURL is unset so callers can surface "not
// configured" distinctly from a transient failure.
func List(proxyURL string) ([]Container, error) {
	if proxyURL == "" {
		return nil, fmt.Errorf("DOCKER_PROXY_URL is not configured")
	}

	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get(strings.TrimRight(proxyURL, "/") + "/containers/json?all=1")
	if err != nil {
		return nil, fmt.Errorf("docker-proxy unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker-proxy returned HTTP %d", resp.StatusCode)
	}

	var raw []rawContainer
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid docker-proxy response: %w", err)
	}

	out := make([]Container, 0, len(raw))
	for _, c := range raw {
		name := strings.TrimPrefix(firstOrEmpty(c.Names), "/")
		health := "none"
		if m := healthPattern.FindStringSubmatch(c.Status); len(m) == 2 {
			health = m[1]
		}
		out = append(out, Container{
			Name:   name,
			Image:  c.Image,
			State:  c.State,
			Status: c.Status,
			Health: health,
		})
	}
	return out, nil
}

func firstOrEmpty(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}
