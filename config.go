package main

import (
	"fmt"
	"strings"
)

type Config struct {
	key string
}

func loadConfig(args []string) (*Config, error) {
	config := &Config{}
	for idx, item := range args {
		i := strings.Index(item, "=")
		if i < 0 {
			return nil, fmt.Errorf("line %d: %w", idx+1, ErrInvalidArgument)
		}
		key := item[:i]
		val := item[i+1:]

		switch key {
		case "key":
			config.key = val
		default:
			fmt.Printf("unknown setting: %v\n", key)
		}
	}
	if config.key == "" {
		return nil, ErrMissingKey
	}

	return config, nil
}
