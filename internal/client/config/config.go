package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		SyncInterval time.Duration `yaml:"syncInterval"`
	} `yaml:"server"`
	App struct {
		WorkDir string `yaml:"workDir"`
	} `yaml:"app"`
}

func (c Config) GetServer() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
