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

func (s *ActionsService) SendShutdown(computerID string, delay int) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "SHUTDOWN", map[string]interface{}{
		"force":   false,
		"timeout": delay,
		"message": fmt.Sprintf("Выключение через %d секунд", delay),
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки команды SHUTDOWN: %w", err)
	}
	return nil
}

func (s *ActionsService) SendBlock(computerID string) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "BLOCK", nil)
	if err != nil {
		return fmt.Errorf("ошибка отправки команды BLOCK: %w", err)
	}
	err = s.actions.BlockComputer(computerID)
	if err != nil {
		fmt.Errorf("ошибка блокировки компьютера в БД: %w", err)
	}
	return nil
}

func (s *ActionsService) SendUnblock(computerID string) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "UNBLOCK", nil)
	if err != nil {
		return fmt.Errorf("ошибка отправки команды UNBLOCK: %w", err)
	}
	err = s.actions.UnblockComputer(computerID)
	if err != nil {
		fmt.Errorf("ошибка разблокировки компьютера в БД: %w", err)
	}
	return nil
}

func (s *ActionsService) SendLockScreen(computerID string) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "LOCK", nil)
	if err != nil {
		return fmt.Errorf("ошибка отправки команды LOCK: %w", err)
	}

	return nil
}

func (s *ActionsService) SendUrl(computerID string, url string) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "URL", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки команды URL: %w", err)
	}
	return nil
}

func (s *ActionsService) SendMessage(computerID string, message string, msgType string) error {
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
	err = zicprotws.SendCommandDirect(tunnelID, "MESSAGE", map[string]interface{}{
		"message": message,
		"type":    msgType,
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки команды MESSAGE: %w", err)
	}
	return nil
}
