package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config содержит настройки бота
type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
	} `yaml:"telegram"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (*Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadLocales загружает строки локализации из YAML-файла.
// LoadLocales загружает строки локализации из YAML-файла.
func LoadLocales(filename string) (map[string]map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла локализации: %w", err)
	}

	var locales map[string]map[string]interface{}
	err = yaml.Unmarshal(data, &locales)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла локализации: %w", err)
	}

	log.Println("Строки локализации успешно загружены")
	return locales, nil
}
