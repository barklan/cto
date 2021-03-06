package caching

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis/v8"
)

var (
	ctx               = context.TODO()
	variableKeySymbol = "$"
)

type Cache interface {
	Set(string, interface{}, time.Duration) error
	Get(string) ([]byte, bool, error)
	SetVar(string, string, interface{}, time.Duration) error
	GetVar(string, string) ([]byte, bool, error)
}

type RedisConnectionData struct {
	Host     string `env:"REDIS_HOST"`
	Password string `env:"REDIS_PASSWORD"`
}

type Redis struct {
	cl *redis.Client
}

func InitRedis(lg *zap.Logger) *Redis {
	cfg := RedisConnectionData{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicln("failed to parse env for redis connection", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", cfg.Host),
		Password: cfg.Password,
		DB:       0, // use default DB
	})
	rs := &Redis{cl: redisClient}

	if err = rs.Set("test", "", 1*time.Minute); err != nil {
		log.Panicln("failed to test test key to redis")
	}
	lg.Info("redis client is ready")

	return rs
}

func (r *Redis) Set(key string, val interface{}, ttl time.Duration) error {
	if ttl < 0 {
		ttl = 0
	}
	err := r.cl.Set(ctx, key, val, ttl).Err()
	return err
}

func (r *Redis) Get(key string) ([]byte, bool, error) {
	val, err := r.cl.Get(ctx, key).Result()
	if err == redis.Nil {
		return []byte{}, false, nil
	} else if err != nil {
		return []byte{}, false, err
	} else {
		return []byte(val), true, nil
	}
}

func (r *Redis) SetVar(namespace, variable string, val interface{}, ttl time.Duration) error {
	fullKey := strings.Join([]string{namespace, variable}, variableKeySymbol)
	return r.Set(fullKey, val, ttl)
}

func (r *Redis) GetVar(namespace, variable string) ([]byte, bool, error) {
	fullKey := strings.Join([]string{namespace, variable}, variableKeySymbol)
	return r.Get(fullKey)
}
