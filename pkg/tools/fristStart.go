package tools

import (
	"backend/pkg/config"
	"backend/pkg/random"
	"fmt"
)

func SetID() error {
	// Генерируем уникальный ID
	id := random.GenerateServerID(10)

	// Сохраняем ID в конфигурацию
	cfg := config.ServiceGet()
	cfg.Server.ID = id

	// Сохраняем изменения
	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
