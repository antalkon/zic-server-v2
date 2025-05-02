package service

import "backend/internal/repository"

type TunelService struct {
	tunel *repository.TunelRepository
}

func NewTunelService(repo *repository.TunelRepository) *TunelService {
	return &TunelService{tunel: repo}
}
