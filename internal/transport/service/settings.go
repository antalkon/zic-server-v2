package service

import (
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"backend/pkg/api"
	"backend/pkg/config"
	"backend/pkg/licenze"
	"fmt"
)

type SettingsService struct {
	settingsRepo *repository.SettingsRepository
}

func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{settingsRepo: repo}
}

type GeneralSettings struct {
	Language            string `json:"language"`
	Timezone            string `json:"timezone"`
	ServerName          string `json:"server_name"`
	ServerType          string `json:"server_type"`
	ServerURL           string `json:"server_url"`
	ServerAddress       string `json:"server_address"`
	ServerPhone         string `json:"server_phone"`
	ServerEmail         string `json:"server_email"`
	ServerContactPerson string `json:"server_contact_person"`
	ServerID            string `json:"server_id"`
	LicenseStatus       string `json:"license_status"`
	LicenseType         string `json:"license_type"`
	LicenseExpiration   string `json:"license_expiration"`
}

type TelegramSettings struct {
	Token           string  `json:"token"`
	Timezone        string  `json:"timezone"`
	AdminIDs        []int64 `json:"admin_ids"`
	TeacherIDs      []int64 `json:"teacher_ids"`
	MessageStart    string  `json:"message_start"`
	MessageHelp     string  `json:"message_help"`
	MessageSettings string  `json:"message_settings"`
}

type LicenseSettings struct {
	Token      string `json:"token"`
	Type       string `json:"type"`
	Expiration string `json:"expiration"`
	AccountID  string `json:"account_id"`
	LastCheck  string `json:"last_check"`
	Status     string `json:"status"`
}
type ApiSettigs struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

func (s *SettingsService) GetGeneralSettings() (*GeneralSettings, error) {
	cfg := config.ServiceGet()

	settings := &GeneralSettings{
		Language:            cfg.Server.Language,
		Timezone:            cfg.Server.Timezone,
		ServerName:          cfg.Server.Name,
		ServerType:          cfg.Server.Type,
		ServerURL:           cfg.Server.URL,
		ServerAddress:       cfg.Server.Address,
		ServerPhone:         cfg.Server.Phone,
		ServerEmail:         cfg.Server.Email,
		ServerContactPerson: cfg.Server.ContactPerson,
		ServerID:            cfg.Server.ID,
		LicenseStatus:       cfg.Licenze.Status,
		LicenseType:         cfg.Licenze.Type,
		LicenseExpiration:   cfg.Licenze.Expiration,
	}

	return settings, nil
}

func (s *SettingsService) UpdateGeneralSettings(req *req.UpdateGeneralSettingsReq, role string) error {
	if role != "admin" {
		return fmt.Errorf("unauthorized")
	}

	cfg := config.ServiceGet()

	if req.Language != nil {
		cfg.Server.Language = *req.Language
	}
	if req.Timezone != nil {
		cfg.Server.Timezone = *req.Timezone
	}
	if req.ServerName != nil {
		cfg.Server.Name = *req.ServerName
	}
	if req.ServerType != nil {
		cfg.Server.Type = *req.ServerType
	}
	if req.ServerURL != nil {
		cfg.Server.URL = *req.ServerURL
	}
	if req.ServerAddress != nil {
		cfg.Server.Address = *req.ServerAddress
	}
	if req.ServerPhone != nil {
		cfg.Server.Phone = *req.ServerPhone
	}
	if req.ServerEmail != nil {
		cfg.Server.Email = *req.ServerEmail
	}
	if req.ServerContactPerson != nil {
		cfg.Server.ContactPerson = *req.ServerContactPerson
	}

	return config.ServiceSaveConfig(cfg)
}

func (s *SettingsService) GetTelegramSettings() (*TelegramSettings, error) {
	cfg := config.TelegramGet()

	settings := &TelegramSettings{
		Token:           cfg.BotSettings.Token,
		Timezone:        cfg.BotSettings.Timezone,
		AdminIDs:        cfg.BotSecurity.AdminIDs,
		TeacherIDs:      cfg.BotSecurity.TeacherIDs,
		MessageStart:    cfg.BotMessages.Start,
		MessageHelp:     cfg.BotMessages.Help,
		MessageSettings: cfg.BotMessages.Settings,
	}

	return settings, nil
}

func (s *SettingsService) UpdateTelegramSettings(req *req.UpdateTelegramSettingsReq, role string) error {
	if role != "admin" {
		return fmt.Errorf("unauthorized")
	}

	cfg := config.TelegramGet()

	if req.Token != nil {
		cfg.BotSettings.Token = *req.Token
	}
	if req.Timezone != nil {
		cfg.BotSettings.Timezone = *req.Timezone
	}
	if req.AdminIDs != nil {
		cfg.BotSecurity.AdminIDs = *req.AdminIDs
	}
	if req.TeacherIDs != nil {
		cfg.BotSecurity.TeacherIDs = *req.TeacherIDs
	}
	if req.MessageStart != nil {
		cfg.BotMessages.Start = *req.MessageStart
	}
	if req.MessageHelp != nil {
		cfg.BotMessages.Help = *req.MessageHelp
	}
	if req.MessageSettings != nil {
		cfg.BotMessages.Settings = *req.MessageSettings
	}

	return config.TelegramSaveConfig(cfg)
}

func (s *SettingsService) GetApiSettings() (*ApiSettigs, error) {
	cfg := config.ServiceGet()

	settings := &ApiSettigs{
		Token: cfg.ZentasAPI.Token,
		Id:    cfg.ZentasAPI.ID,
	}

	return settings, nil
}

func (s *SettingsService) UpdateApiSettings(req *req.UpdateApiSettingsReq, role string) error {
	if role != "admin" {
		return fmt.Errorf("unauthorized")
	}
	if err := api.ActivateApi(req.Id, req.Token); err != nil {
		return fmt.Errorf("activation error: %w", err)
	}
	cfg := config.ServiceGet()
	cfg.ZentasAPI.Token = req.Token
	cfg.ZentasAPI.ID = req.Id

	return config.ServiceSaveConfig(cfg)
}

func (s *SettingsService) GetLicenseSettings() (*LicenseSettings, error) {
	cfg := config.ServiceGet()

	settings := &LicenseSettings{
		Token:      cfg.Licenze.Token,
		Type:       cfg.Licenze.Type,
		Expiration: cfg.Licenze.Expiration,
		AccountID:  cfg.Licenze.AccountID,
		LastCheck:  cfg.Licenze.LastCheck,
		Status:     cfg.Licenze.Status,
	}

	return settings, nil
}

func (s *SettingsService) UpdateLicenseSettings(req *req.UpdateLicenseSettingsReq, role string) error {
	if role != "admin" {
		return fmt.Errorf("unauthorized")
	}

	_, err := licenze.ActivateLicenze(req.Token)
	if err != nil {
		return fmt.Errorf("activation error: %w", err)
	}

	return nil
}
