package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const telegramConfigPath = "config/yaml/telegram.yaml"

var (
	telegramCfg  *TelegramConfig
	telegramOnce sync.Once
)

type TelegramConfig struct {
	BotSettings TelegramBotSettings `yaml:"bot-settings"`
	BotSecurity TelegramBotSecurity `yaml:"bot-seccurity"`
	BotMessages TelegramBotMessages `yaml:"bot-messages"`
}

type TelegramBotSettings struct {
	Token    string `yaml:"token"`
	Timezone string `yaml:"timezone"`
}

type TelegramBotSecurity struct {
	AdminIDs   []int64 `yaml:"admin-ids"`
	TeacherIDs []int64 `yaml:"teacher-ids"`
}

type TelegramBotMessages struct {
	Start    string `yaml:"start"`
	Help     string `yaml:"help"`
	Settings string `yaml:"settings"`
}

// TelegramInit — инициализирует и загружает конфиг
func TelegramInit() {
	telegramOnce.Do(func() {
		var err error
		telegramCfg, err = telegramLoad()
		if err != nil {
			panic("❌ Failed to load telegram config: " + err.Error())
		}
	})
}

// TelegramGet — глобальный доступ к Telegram конфигу
func TelegramGet() *TelegramConfig {
	if telegramCfg == nil {
		panic("Telegram config not initialized — call config.TelegramInit() first")
	}
	return telegramCfg
}

// TelegramSave — сохраняет конфиг
func TelegramSaveConfig(cfg *TelegramConfig) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("configs", 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(telegramConfigPath, out, 0644)
}

// telegramLoad — загружает и делает миграцию конфига
func telegramLoad() (*TelegramConfig, error) {
	if _, err := os.Stat(telegramConfigPath); os.IsNotExist(err) {
		fmt.Println("⚙️ Telegram config not found. Creating a template...")
		cfg := &TelegramConfig{}
		if err := TelegramSaveConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := ioutil.ReadFile(telegramConfigPath)
	if err != nil {
		return nil, err
	}

	var parsed TelegramConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	finalYaml, err := yaml.Marshal(&parsed)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(telegramConfigPath, finalYaml, 0644); err != nil {
		return nil, err
	}

	return &parsed, nil
}
