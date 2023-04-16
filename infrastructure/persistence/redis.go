package persistence_infrastructure

import (
	"github.com/redis/go-redis/v9"
)

func NewRedis(conf redis.Options) *redis.Client {
	return redis.NewClient(&conf)
}
