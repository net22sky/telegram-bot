package config

import (
    "fmt"
    "log"
    "os"

    "gopkg.in/yaml.v3"
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

// LoadLocales загружает строки локализации из файла
func LoadLocales(filename string) (map[string]map[string]string, error) {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return nil, fmt.Errorf("locales file not found: %s", filename)
    }

    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var locales map[string]map[string]string
    if err := yaml.Unmarshal(data, &locales); err != nil {
        return nil, err
    }

    return locales, nil
}