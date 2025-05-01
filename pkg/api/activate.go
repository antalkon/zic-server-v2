package api

import (
	glob_config "backend/config"
	"backend/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type bodyActivate struct {
	ID            string `json:"id"`
	Token         string `json:"token"`
	ServerID      string `json:"server_id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	URL           string `json:"url"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	ContactPerson string `json:"contact_person"`
	LicenseToken  string `json:"license_token"`
}

func ActivateApi(id, token string) error {
	cfg := config.ServiceGet()
	url := glob_config.ZicApiUrl + "/activate"

	body := bodyActivate{
		ID:            id,
		Token:         token,
		ServerID:      cfg.Server.ID,
		Name:          cfg.Server.Name,
		Type:          cfg.Server.Type,
		URL:           cfg.Server.URL,
		Address:       cfg.Server.Address,
		Phone:         cfg.Server.Phone,
		Email:         cfg.Server.Email,
		ContactPerson: cfg.Server.ContactPerson,
		LicenseToken:  cfg.Licenze.Token,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("activation failed: %s", respBody)
	}

	fmt.Println("✅ API успешно активирован:", string(respBody))
	return nil
}
