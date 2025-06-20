package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var redisClient = redis.NewClient

type Cache struct {
	*redis.Client
}

// NewCache创建一个Redis客户端实例
func NewCache(config *viper.Viper) (*Cache, error) {
	// 创建redis  client 客户端实例
	client := redisClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.GetString("redis.host"), config.GetInt("redis.port")),
		Password:     config.GetString("redis.password"),
		DB:           config.GetInt("redis.db"),
		PoolSize:     config.GetInt("redis.poolSize"),
		MinIdleConns: config.GetInt("redis.minIdleConns"),
		MaxIdleConns: config.GetInt("redis.maxIdelConns"),
	})
	//  测试连接
	// 创建带有超时机制的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {

		return nil, fmt.Errorf("连接Redis失败 %v", err)
	}

	return &Cache{client}, nil
}

// Close 关闭管道资源
func (c *Cache) Close() error {

	if err := c.Client.Close(); err != nil {
		return fmt.Errorf("关闭redis资源失败 %v", err)
	}
	return nil
}
