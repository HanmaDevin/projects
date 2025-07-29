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
	Prompt string `yaml:"prompt"`
	Model  string `yaml:"model"`
	Stream bool   `yaml:"stream"`
}

func ReadConfig() *ollama.OllamaModel {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		WriteConfig(Config{
			Prompt: "What is the meaning of life?",
			Model:  "",
			Stream: false,
		})
		return nil
	}
	err = yaml.Unmarshal(data, &cfg)

	return parseConfig(cfg)
}

func WriteConfig(cfg Config) error {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func parseConfig(cfg Config) *ollama.OllamaModel {
	Body := ollama.NewOllamaModel()
	Body.Model = cfg.Model
	Body.Prompt = cfg.Prompt
	Body.Stream = cfg.Stream
	return Body
}
