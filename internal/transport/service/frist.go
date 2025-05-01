package service

import (
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"backend/pkg/api"
	"backend/pkg/config"
	"backend/pkg/licenze"
	"fmt"
)

type FristService struct {
	fristRepo *repository.FristRepository
}

func NewFristService(repo *repository.FristRepository) *FristService {
	return &FristService{fristRepo: repo}
}

func (s *FristService) ActivateLicenze(token string) error {
	_, err := licenze.ActivateLicenze(token)
	if err != nil {
		return fmt.Errorf("failed to activate license: %w", err)
	}
	return nil
}

func (s *FristService) OrgFormation(req *req.FristStartForm) error {
	cfg := config.ServiceGet()

	cfg.Server.Name = req.Name
	cfg.Server.Type = req.Type
	cfg.Server.URL = req.Url
	cfg.Server.Address = req.Address
	cfg.Server.Phone = req.Phone
	cfg.Server.Email = req.Email
	cfg.Server.ContactPerson = req.ContactPerson

	// Сохраняем изменения
	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (s *FristService) ApiFormation(req *req.FristStartAPI) error {
	id := req.Id
	token := req.Token

	if err := api.ActivateApi(id, token); err != nil {
		return fmt.Errorf("failed to activate API: %w", err)
	}

	// Сохраняем изменения
	cfg := config.ServiceGet()
	cfg.ZentasAPI.ID = id
	cfg.ZentasAPI.Token = token
	cfg.Server.Frist = false
	if err := config.ServiceSaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}
