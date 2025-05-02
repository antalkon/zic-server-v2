package tunel_service

import (
	"backend/internal/repository"
	wsmodels "backend/internal/ws_models"
	"backend/pkg/cache"
	tokenjwt "backend/pkg/token_jwt"
	"fmt"

	"github.com/google/uuid"
)

type TunelService struct {
	tunel *repository.TunelRepository
	cache *cache.RedisClient
}

func NewTunelService(repo *repository.TunelRepository, cache *cache.RedisClient) *TunelService {
	return &TunelService{tunel: repo, cache: cache}
}

// TODO: Реализовать функционал Secret!!!
func (s *TunelService) GetTunnelByID(ws *wsmodels.InitPayload) error {
	id, err := uuid.Parse(ws.ComputerID)
	if err != nil {
		return fmt.Errorf("invalid computer ID format: %w", err)
	}

	pc, err := s.tunel.GetPcById(id)
	if err != nil {
		return fmt.Errorf("failed to get computer: %w", err)
	}

	jwt, err := tokenjwt.DecodePcJWT(ws.JwtToken)
	if err != nil {
		return fmt.Errorf("failed to decode JWT: %w", err)
	}

	jwtUUID, err := uuid.Parse(jwt.ID)
	if err != nil {
		return fmt.Errorf("invalid JWT ID format: %w", err)
	}

	if pc.ID != jwtUUID {
		return fmt.Errorf("computer ID does not match JWT")
	}

	if pc.ClientVersion != ws.ClientVersion {
		if err := s.tunel.UpdateClientVersion(pc.ID, ws.ClientVersion); err != nil {
			return fmt.Errorf("failed to update client version: %w", err)
		}
	}

	if pc.OS != ws.OS {
		if err := s.tunel.UpdateOS(pc.ID, ws.OS); err != nil {
			return fmt.Errorf("failed to update OS: %w", err)
		}
	}

	if pc.PublicIP != ws.PublicIP {
		if err := s.tunel.UpdatePublicIP(pc.ID, ws.PublicIP); err != nil {
			return fmt.Errorf("failed to update public IP: %w", err)
		}
	}

	if pc.LocalIP != ws.LocalIP {
		if err := s.tunel.UpdateLocalIP(pc.ID, ws.LocalIP); err != nil {
			return fmt.Errorf("failed to update local IP: %w", err)
		}
	}
	if err := s.tunel.UpdateLastActivity(pc.ID); err != nil {
		return fmt.Errorf("failed to update last activity: %w", err)
	}

	return nil
}
