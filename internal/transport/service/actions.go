package service

import (
	"backend/internal/repository"
	zicprotws "backend/internal/transport/zicprot_ws"
	"backend/pkg/cache"
	"encoding/json"
	"fmt"
)

type ActionsService struct {
	actions *repository.ActionsRepository
	cache   *cache.RedisClient
}

func NewActionsService(repo *repository.ActionsRepository, cache *cache.RedisClient) *ActionsService {
	return &ActionsService{actions: repo, cache: cache}
}

func (s *ActionsService) SendReboot(computerID string, delay int) error {
	key := "pc:" + computerID

	raw, err := s.cache.Get(key)
	if err != nil {
		return fmt.Errorf("не удалось получить данные по ключу %s: %w", key, err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return fmt.Errorf("ошибка парсинга JSON из Redis: %w", err)
	}

	tunnelID, ok := data["tunnel_id"].(string)
	if !ok {
		return fmt.Errorf("поле tunnel_id отсутствует или неправильного типа")
	}

	// 💡 Тут передаём tunnelID вместо computerID!
	err = zicprotws.SendCommandDirect(tunnelID, "REBOOT", map[string]interface{}{
		"force":   false,
		"timeout": delay,
		"message": fmt.Sprintf("Перезагрузка через %d секунд", delay),
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки команды REBOOT: %w", err)
	}
	return nil
}
