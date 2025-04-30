package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const serviceConfigPath = "config/service.yaml"

var (
	serviceCfg  *ServiceConfig
	serviceOnce sync.Once
)

type ServiceConfig struct {
	Licenze   ServiceLicenzeConfig `yaml:"licenze"`
	Server    ServiceServerConfig  `yaml:"server"`
	ZentasAPI ServiceZentasAPI     `yaml:"zentas-api"`
}

type ServiceLicenzeConfig struct {
	Token      string `yaml:"token"`
	Type       string `yaml:"type"`
	Expiration string `yaml:"expiration"`
	AccountID  string `yaml:"accound-id"`
	LastCheck  string `yaml:"last-check"`
	Secret     string `yaml:"secret"`
	Status     string `yaml:"status"`
}

type ServiceServerConfig struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	URL           string `yaml:"url"`
	Address       string `yaml:"address"`
	Phone         string `yaml:"phone"`
	Email         string `yaml:"email"`
	ContactPerson string `yaml:"contact-person"`
	ID            string `yaml:"id"`
}

type ServiceZentasAPI struct {
	Token string `yaml:"token"`
	ID    string `yaml:"id"`
}

// ServiceInit — загружает или создаёт конфиг
func ServiceInit() {
	serviceOnce.Do(func() {
		var err error
		serviceCfg, err = serviceLoad()
		if err != nil {
			panic("❌ Failed to load service config: " + err.Error())
		}
	})
}

// ServiceGet — безопасный глобальный доступ к конфигу
func ServiceGet() *ServiceConfig {
	if serviceCfg == nil {
		panic("Service config not initialized — call config.ServiceInit() first")
	}
	return serviceCfg
}

// ServiceSave — сохраняет YAML
func ServiceSaveConfig(cfg *ServiceConfig) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("configs", 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(serviceConfigPath, out, 0644)
}

// serviceLoad — читает и дополняет YAML
func serviceLoad() (*ServiceConfig, error) {
	if _, err := os.Stat(serviceConfigPath); os.IsNotExist(err) {
		fmt.Println("⚙️ Service config not found. Creating a template...")
		cfg := &ServiceConfig{}
		if err := ServiceSaveConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := ioutil.ReadFile(serviceConfigPath)
	if err != nil {
		return nil, err
	}

	var parsed ServiceConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	// Миграция: добавление недостающих полей
	finalYaml, err := yaml.Marshal(&parsed)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(serviceConfigPath, finalYaml, 0644); err != nil {
		return nil, err
	}

	return &parsed, nil
}
