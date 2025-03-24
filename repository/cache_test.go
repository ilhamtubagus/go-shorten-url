package repository

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/ilhamtubagus/go-shorten-url/constants"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func setupRedisCache() (*RedisCache[TestStruct], *miniredis.Miniredis, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := NewRedisCache[TestStruct](client)
	return cache, mr, nil
}

func TestRedisCache_Put(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()
	testData := TestStruct{ID: 1, Name: "Test"}

	err = cache.Put(ctx, "test_key", testData, 60)
	assert.NoError(t, err)

	// Verify data was stored correctly
	val, err := mr.Get("test_key")
	assert.NoError(t, err)

	var storedData TestStruct
	err = json.Unmarshal([]byte(val), &storedData)
	assert.NoError(t, err)
	assert.Equal(t, testData, storedData)
}

func TestRedisCache_Get(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()
	testData := TestStruct{ID: 1, Name: "Test"}

	// Store test data
	jsonData, _ := json.Marshal(testData)
	mr.Set("test_key", string(jsonData))

	// Test Get
	result, err := cache.Get(ctx, "test_key")
	assert.NoError(t, err)
	assert.Equal(t, testData, result)

	// Test Get with non-existent key
	_, err = cache.Get(ctx, "non_existent_key")
	assert.ErrorIs(t, err, constants.ErrorCacheNotFound)
}

func TestRedisCache_IsExist(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()

	// Test non-existent key
	exists, err := cache.IsExist(ctx, "test_key")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Set a key
	mr.Set("test_key", "value")

	// Test existing key
	exists, err = cache.IsExist(ctx, "test_key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisCache_Delete(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()

	// Set a key
	mr.Set("test_key", "value")

	// Delete the key
	err = cache.Delete(ctx, "test_key")
	assert.NoError(t, err)

	// Verify key was deleted
	exists, err := cache.IsExist(ctx, "test_key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisCache_Flush(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()

	// Set multiple keys
	mr.Set("key1", "value1")
	mr.Set("key2", "value2")

	// Flush all keys
	err = cache.Flush(ctx)
	assert.NoError(t, err)

	// Verify all keys were deleted
	keys := mr.Keys()
	assert.NoError(t, err)
	assert.Empty(t, keys)
}

func TestRedisCache_PutWithTTL(t *testing.T) {
	cache, mr, err := setupRedisCache()
	assert.NoError(t, err)
	defer mr.Close()

	ctx := context.Background()
	testData := TestStruct{ID: 1, Name: "Test"}

	err = cache.Put(ctx, "test_key", testData, 1) // 1 second TTL
	assert.NoError(t, err)

	// Verify data was stored correctly
	val, err := mr.Get("test_key")
	assert.NoError(t, err)

	var storedData TestStruct
	err = json.Unmarshal([]byte(val), &storedData)
	assert.NoError(t, err)
	assert.Equal(t, testData, storedData)

	// Verify the key exists
	exists, err := cache.IsExist(ctx, "test_key")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Fast-forward time in miniredis
	mr.FastForward(2 * time.Second)

	// Verify the key has expired
	exists, err = cache.IsExist(ctx, "test_key")
	assert.NoError(t, err)
	assert.False(t, exists)
}
