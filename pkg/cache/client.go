package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/umutcomlekci/automated-messaging-system/internal/config"
)

type Client struct {
	redisClient *redis.Client
}

func NewCacheClient() (*Client, error) {
	opt, err := redis.ParseURL(config.GetCacheConnectionString())
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Client{
		redisClient: redisClient,
	}, nil
}

func (r *Client) SetStruct(key string, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = r.redisClient.Set(context.Background(), key, jsonValue, 0).Result()
	return err
}

func (r *Client) GetStruct(key string, pointer any) error {
	value, err := r.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(value), &pointer)
	return err
}
