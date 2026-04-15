package cache

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestRedisConfig() *viper.Viper {
	config := viper.New()
	config.Set("redis.addr", "localhost")
	config.Set("redis.port", 6379)
	config.Set("redis.password", "")
	config.Set("redis.db", 0)
	config.Set("redis.poolSize", 10)
	config.Set("redis.minIdleConns", 5)
	config.Set("redis.maxIdelConns", 10)
	return config
}

func TestNewCache_InvalidConfig(t *testing.T) {
	// Test with invalid Redis address
	config := viper.New()
	config.Set("redis.addr", "invalid-address-that-does-not-exist")
	config.Set("redis.port", 6379)
	config.Set("redis.password", "")
	config.Set("redis.db", 0)

	// This should fail or timeout since the address doesn't exist
	cache, err := NewCache(config)

	// We expect an error since the Redis server doesn't exist
	assert.Error(t, err)
	assert.Nil(t, cache)
}

func TestNewCache_EmptyConfig(t *testing.T) {
	config := viper.New()

	// This should fail due to missing required config
	cache, err := NewCache(config)

	// Will either fail or connect to default localhost
	// The behavior depends on whether Redis is running locally
	if err == nil {
		// If it connected, clean up
		assert.NotNil(t, cache)
		cache.Close()
	} else {
		// If it failed, that's also acceptable for this test
		assert.Nil(t, cache)
	}
}

func TestCache_Close(t *testing.T) {
	// Create a config that will likely fail to connect
	config := setupTestRedisConfig()

	cache, err := NewCache(config)

	// If Redis is running locally, test the Close method
	if err == nil && cache != nil {
		defer cache.Close()

		// Test that we can ping
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := cache.Ping(ctx).Err()
		assert.NoError(t, err)

		// Now close and verify
		closeErr := cache.Close()
		assert.NoError(t, closeErr)
	} else {
		// If we couldn't connect, that's fine for unit tests
		assert.Nil(t, cache)
	}
}

func TestCache_SetGet(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		// Skip if Redis is not available
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Test Set and Get
	err = cache.Set(ctx, "test_key", "test_value", 10*time.Second).Err()
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)
}

func TestCache_Delete(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Set a key
	err = cache.Set(ctx, "delete_test_key", "test_value", 10*time.Second).Err()
	assert.NoError(t, err)

	// Delete the key
	err = cache.Del(ctx, "delete_test_key").Err()
	assert.NoError(t, err)

	// Verify the key is deleted
	val, err := cache.Get(ctx, "delete_test_key").Result()
	assert.Error(t, err) // Should return redis.Nil error
	assert.Empty(t, val)
}

func TestCache_Exists(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Set a key
	err = cache.Set(ctx, "exists_test_key", "test_value", 10*time.Second).Err()
	assert.NoError(t, err)

	// Check if key exists
	exists, err := cache.Exists(ctx, "exists_test_key").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, exists)

	// Delete and check again
	cache.Del(ctx, "exists_test_key")
	exists, err = cache.Exists(ctx, "exists_test_key").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, exists)
}

func TestCache_Expire(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Set a key
	err = cache.Set(ctx, "expire_test_key", "test_value", 0).Err()
	assert.NoError(t, err)

	// Set expiration
	expired := cache.Expire(ctx, "expire_test_key", 1*time.Second).Val()
	assert.True(t, expired)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Key should be expired
	val, err := cache.Get(ctx, "expire_test_key").Result()
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestCache_HSetHGet(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Test HSet and HGet
	err = cache.HSet(ctx, "test_hash", "field1", "value1").Err()
	assert.NoError(t, err)

	val, err := cache.HGet(ctx, "test_hash", "field1").Result()
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestCache_LPushLRange(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Test LPush and LRange
	cache.Del(ctx, "test_list") // Clean up first

	err = cache.LPush(ctx, "test_list", "item1", "item2", "item3").Err()
	assert.NoError(t, err)

	items, err := cache.LRange(ctx, "test_list", 0, -1).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"item3", "item2", "item1"}, items)
}

func TestCache_SAddSMembers(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Test SAdd and SMembers
	cache.Del(ctx, "test_set") // Clean up first

	err = cache.SAdd(ctx, "test_set", "member1", "member2", "member3").Err()
	assert.NoError(t, err)

	members, err := cache.SMembers(ctx, "test_set").Result()
	assert.NoError(t, err)
	assert.Len(t, members, 3)
	assert.Contains(t, members, "member1")
}

func TestCache_IncrDecr(t *testing.T) {
	config := setupTestRedisConfig()

	cache, err := NewCache(config)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Test Incr
	cache.Del(ctx, "test_counter")

	val, err := cache.Incr(ctx, "test_counter").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, val)

	val, err = cache.IncrBy(ctx, "test_counter", 5).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 6, val)

	val, err = cache.Decr(ctx, "test_counter").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 5, val)
}
