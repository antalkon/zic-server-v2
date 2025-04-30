package service

import (
	"backend/internal/repository"
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
