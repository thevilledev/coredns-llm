package llm

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

const (
	envAPIKey = "COREDNS_LLM_API_KEY"
)

// init registers the plugin.
func init() {
	plugin.Register("llm", setup)
}

// setup configures the llm plugin.
func setup(c *caddy.Controller) error {
	config, err := parse(c)
	if err != nil {
		return plugin.Error("llm", err)
	}

	apiKey := os.Getenv(envAPIKey)
	if apiKey == "" {
		return plugin.Error("llm", fmt.Errorf("environment variable %s is not set", envAPIKey))
	}

	client := &OpenAIClient{
		EndpointURL: config.EndpointURL,
		Model:       config.Model,
		APIKey:      apiKey,
		HTTPClient:  &http.Client{Timeout: time.Duration(config.TimeoutSeconds) * time.Second},
	}

	h := &Handler{
		ChunkSize: config.ChunkSize,
		Client:    client,
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		h.Next = next
		return h
	})

	return nil
}

// parse reads llm configuration from the Corefile.
func parse(c *caddy.Controller) (Config, error) {
	cfg := defaultConfig()

	for c.Next() { // llm
		for c.NextBlock() {
			switch c.Val() {
			case "endpoint":
				if !c.NextArg() {
					return cfg, c.ArgErr()
				}
				cfg.EndpointURL = c.Val()
			case "model":
				if !c.NextArg() {
					return cfg, c.ArgErr()
				}
				cfg.Model = c.Val()
			case "chunk_size":
				if !c.NextArg() {
					return cfg, c.ArgErr()
				}
				s, err := parsePositiveInt(c.Val())
				if err != nil {
					return cfg, c.Errf("chunk_size: %v", err)
				}
				cfg.ChunkSize = s
			case "timeout":
				if !c.NextArg() {
					return cfg, c.ArgErr()
				}
				s, err := parsePositiveInt(c.Val())
				if err != nil {
					return cfg, c.Errf("timeout: %v", err)
				}
				cfg.TimeoutSeconds = s
			default:
				return cfg, c.Errf("unknown property %q", c.Val())
			}
		}
	}

	if cfg.EndpointURL == "" {
		return cfg, fmt.Errorf("endpoint is required")
	}
	if cfg.Model == "" {
		return cfg, fmt.Errorf("model is required")
	}

	return cfg, nil
}
