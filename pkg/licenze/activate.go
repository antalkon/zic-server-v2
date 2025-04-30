package licenze

import (
	glob_config "backend/config"
	"backend/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ActivateLicenze(token string) (bool, error) {
	if glob_config.Licenze {
		return false, fmt.Errorf("лицензия активна")
	}
	if err := activateReq(token); err != nil {
		return false, err
	}
	fmt.Println("✅ Лицензия успешно активирована")
	return true, nil

}

func activateReq(token string) error {
	cfg := config.ServiceGet()
	id := cfg.Server.ID
	url := glob_config.ZicApiUrl + "/activate"

	payload := map[string]string{
		"token":      token,
		"service_id": id,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("activation error: %s", body)
	}

	var result struct {
		Expiration string `json:"expiration"`
		NewSecret  string `json:"new_secret"`
		AccountID  string `json:"account-id"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse response error: %w", err)
	}

	// Обновляем конфиг
	cfg.Licenze.Expiration = result.Expiration
	cfg.Licenze.Secret = result.NewSecret
	cfg.Licenze.AccountID = result.AccountID

	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✅ Лицензия успешно активирована")
	return nil
}
