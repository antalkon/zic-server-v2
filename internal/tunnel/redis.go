package tunnel

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{Client: client}
}

func (r *RedisClient) SaveTunnel(tunnel *Tunnel) error {
	json, err := json.Marshal(tunnel)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Client.Set(ctx, tunnel.ID, json, 2*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetTunnel(id string) (*Tunnel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonStr, err := r.Client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	tunnel := &Tunnel{}
	err = json.Unmarshal([]byte(jsonStr), tunnel)
	if err != nil {
		return nil, err
	}

	return tunnel, nil
}

func (r *RedisClient) DeleteTunnel(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.Client.Del(ctx, id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) PublishToTunnel(tunnelID string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.Client.Publish(ctx, tunnelID, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) UpdateTunnelTTL(tunnelID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.Client.Expire(ctx, tunnelID, 2*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

type Tunnel struct {
	ID         string `json:"id"`
	ComputerID string `json:"computer_id"`
	CreatedAt  int64  `json:"created_at"`
}
