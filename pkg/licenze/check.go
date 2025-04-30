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
	now := time.Now()

	// ✅ Если уже была проверка сегодня и статус "active" — пропускаем
	if last, err := time.Parse("2006-01-02", cfg.Licenze.LastCheck); err == nil {
		if last.Format("2006-01-02") == now.Format("2006-01-02") && cfg.Licenze.Status == "active" {
			fmt.Println("✅ Лицензия активна и уже проверялась сегодня")
			glob_config.Licenze = true
			return true
		}
	}

	// Проверяем срок действия
	exp, err := time.Parse("2006-01-02", cfg.Licenze.Expiration)
	if err != nil || now.After(exp) || cfg.Licenze.Status != "active" {
		fmt.Println("⚠️ Лицензия невалидна или истекла. Проверяем заново...")
		if err := checkReq(); err != nil {
			fmt.Println("❌ Ошибка проверки лицензии:", err)
			glob_config.Licenze = false
			return false
		}
	}

	// Если проверка не была сегодня, делаем её
	if err := checkReq(); err != nil {
		fmt.Println("❌ Ошибка проверки лицензии:", err)
		glob_config.Licenze = false
		return false
	}

	fmt.Println("✅ Лицензия успешно проверена")
	glob_config.Licenze = true
	return true
}
func checkReq() error {
	cfg := config.ServiceGet()
	url := glob_config.ZicApiUrl + "/check"

	payload := map[string]string{
		"token":  cfg.Licenze.Token,
		"secret": cfg.Licenze.Secret,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("invalid response: %s", bodyBytes)
	}

	var res struct {
		Expiration string `json:"expiration"` // 🟢 Исправлено имя
		NewSecret  string `json:"new_secret"`
		AccountID  string `json:"account_id"`
		Token      string `json:"token"` // опционально, но вдруг вернёт
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	// Обновляем конфиг
	cfg.Licenze.Expiration = res.Expiration
	cfg.Licenze.Secret = res.NewSecret
	cfg.Licenze.AccountID = res.AccountID
	cfg.Licenze.LastCheck = time.Now().Format("2006-01-02")
	cfg.Licenze.Status = "active"
	if res.Token != "" {
		cfg.Licenze.Token = res.Token
	}

	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("save config error: %w", err)
	}

	fmt.Println("✅ Конфигурация лицензии обновлена")
	return nil
}
