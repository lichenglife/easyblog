package cache

import (
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setMockViperConfig() *viper.Viper {

	viper := viper.New()
	// viper := &mockViperConfig{
	// 	viper: viper.New(),
	// }
	viper.Set("cache", "redis")
	viper.Set("redis.host", "localhost")
	viper.Set("redis.port", 6379)
	viper.Set("redis.password", "")
	viper.Set("redis.db", 0)
	viper.Set("redis.poolSize", 10)
	viper.Set("redis.minIdleConns", 10)
	viper.Set("redis.maxIdelConns", 10)

	return viper
}

func TestNewCache_Success(t *testing.T) {
	// 1、设置viper
	config := setMockViperConfig()

	// 2、创建redis mock
	client, mock := redismock.NewClientMock()

	// 3、设置mock 期望
	mock.ExpectPing().SetVal("PONG")

	// 4、替换真实的redis.NewClient
	originalNewClient := redisClient
	redisClient = func(opt *redis.Options) *redis.Client {
		return client
	}
	defer func() {
		redisClient = originalNewClient
	}()
	// 5、调用 NewCache
	cache, err := NewCache(config)
	// 6、断言
	assert.Nil(t, err)
	assert.NotNil(t, cache)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试关闭
	err = cache.Close()

	assert.NoError(t, err)

}

func TestNewCache_PING_Failed(t *testing.T) {
	// 1、设置viper
	config := setMockViperConfig()

	// 2、创建redis mock
	client, mock := redismock.NewClientMock()

	// 3、设置mock 期望
	mock.ExpectPing().SetErr(errors.New("ping failed"))

	// 4、替换真实的redis.NewClient
	originalNewClient := redisClient
	redisClient = func(opt *redis.Options) *redis.Client {
		return client
	}
	defer func() {
		redisClient = originalNewClient
	}()
	// 5、调用 NewCache
	cache, err := NewCache(config)
	// 6、断言
	assert.Nil(t, cache)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "ping failed")
}
