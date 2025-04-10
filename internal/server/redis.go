package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kevinhartarto/tasker/internal/logger"
	"github.com/kevinhartarto/tasker/internal/utils"
	"github.com/redis/go-redis/v9"
)

var (
	ctx       = context.Background()
	serverLog = logger.GetLogger()
)

func StartRedis() *redis.Client {
	redisUrl := utils.GetEnvOrDefault("REDIS_URL", "redis")
	redisPort := utils.GetEnvOrDefault("REDIS_PORT", "6379")
	redisAddr := fmt.Sprintf("%v:%v", redisUrl, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		serverLog.Info("Could not connect to redis: " + err.Error())
		os.Exit(1)
	}
	serverLog.Info("Connected to redis: " + pong)

	return redisClient
}
