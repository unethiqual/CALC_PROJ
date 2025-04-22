package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TimeAdditionMs        int `json:"time_addition_ms"`
	TimeSubtractionMs     int `json:"time_subtraction_ms"`
	TimeMultiplicationsMs int `json:"time_multiplications_ms"`
	TimeDivisionsMs       int `json:"time_divisions_ms"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
