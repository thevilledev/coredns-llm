package llm

import (
	"fmt"
)

// Config holds Corefile configuration for the plugin.
type Config struct {
	EndpointURL    string
	Model          string
	ChunkSize      int
	TimeoutSeconds int
}

func defaultConfig() Config {
	return Config{
		ChunkSize:      200,
		TimeoutSeconds: 15,
	}
}

// parsePositiveInt converts a string to a positive int.
func parsePositiveInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("must be a positive integer")
	}
	return n, nil
}
