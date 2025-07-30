package config

import (
	"os"
	"path/filepath"

	"github.com/HanmaDevin/schlama/ollama"
	"gopkg.in/yaml.v3"
)

var home, _ = os.UserHomeDir()
var config_Path string = filepath.Dir(home + "/.config/schlama/")
var filename string = config_Path + "/config.yaml"

type Config struct {
	Model string `yaml:"model"`
}

func ReadConfig() *ollama.OllamaModel {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		WriteConfig(Config{
			Model: "",
		})
		return nil
	}
	// ignore errors, there should'nt be any
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil
	}

	return parseConfig(cfg)
}

func WriteConfig(cfg Config) error {
	// Ensure the config directory exists
	if err := os.MkdirAll(config_Path, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func parseConfig(cfg Config) *ollama.OllamaModel {
	Body := ollama.NewOllamaModel()
	Body.Model = cfg.Model
	Body.Prompt = ""
	Body.Stream = false
	return Body
}

// UpdateModel updates just the model in the configuration
func UpdateModel(modelName string) error {
	cfg := Config{
		Model: modelName,
	}
	return WriteConfig(cfg)
}
