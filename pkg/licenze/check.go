package licenze

import (
	glob_config "backend/config"
	"backend/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func CheckLicenze() bool {
	cfg := config.ServiceGet()
	status := false
	// Парсим дату истечения лицензии
	licenzeExpiration, err := time.Parse("2006-01-02", cfg.Licenze.Expiration)
	if err != nil {
		status = false
	}

	currentDay := time.Now()

	// Проверяем срок действия
	if currentDay.After(licenzeExpiration) {
		status = false
	}

	// Проверяем статус
	if cfg.Licenze.Status != "active" {
		status = false
	}

	// Проверка, была ли уже выполнена проверка сегодня
	lastCheck, err := time.Parse("2006-01-02", cfg.Licenze.LastCheck)
	if err == nil && currentDay.Format("2006-01-02") == lastCheck.Format("2006-01-02") {
		fmt.Println("✅ Лицензия активна и проверена")
		glob_config.Licenze = true
		return true
	}

	if !status {
		if err := checkReq(); err != nil {
			fmt.Println("❌ Ошибка проверки лицензии:", err)
			return false
		}
	}
	fmt.Println("✅ Лицензия активна и проверена")
	glob_config.Licenze = true
	return true
}

func checkReq() error {
	lic := config.ServiceGet().Licenze
	url := glob_config.ZicApiUrl + "/check"
	payload := map[string]string{
		"token":  lic.Token,
		"secret": lic.Secret,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s", content)
	}

	var data struct {
		Expiration string `json:"expiratiom"`
		NewSecret  string `json:"new secret"`
		AccountID  string `json:"account-id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Обновляем поля в конфиге
	cfg := config.ServiceGet()
	cfg.Licenze.Expiration = data.Expiration
	cfg.Licenze.Secret = data.NewSecret
	cfg.Licenze.AccountID = data.AccountID
	cfg.Licenze.LastCheck = time.Now().Format("2006-01-02")
	cfg.Licenze.Status = "active"

	// Сохраняем обновлённый конфиг
	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✅ Лицензия успешно проверена и обновлена")
	return nil
}
