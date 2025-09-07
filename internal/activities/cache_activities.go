package activities

import (
	"github.com/umutcomlekci/automated-messaging-system/pkg/cache"
)

type CacheActivities struct {
	cacheClient *cache.Client
}

func NewCacheActivities(cacheClient *cache.Client) *CacheActivities {
	return &CacheActivities{
		cacheClient: cacheClient,
	}
}

func (a *CacheActivities) SetStruct(key string, value *cache.SentMessageCache) error {
	return a.cacheClient.SetStruct(key, value)
}

func (a *CacheActivities) GetStruct(key string) (*cache.SentMessageCache, error) {
	var message *cache.SentMessageCache
	err := a.cacheClient.GetStruct(key, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
