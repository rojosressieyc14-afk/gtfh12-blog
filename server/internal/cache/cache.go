package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// GetOrSet tries to get key from Redis; on miss, calls dbQuery and caches the result.
func GetOrSet(ctx context.Context, key string, ttl time.Duration, dest interface{}, dbQuery func() error) error {
	if RDB == nil {
		return dbQuery()
	}

	val, err := RDB.Get(ctx, key).Result()
	if err == nil {
		return json.Unmarshal([]byte(val), dest)
	}

	if err := dbQuery(); err != nil {
		return err
	}

	data, err := json.Marshal(dest)
	if err != nil {
		return nil
	}
	RDB.Set(ctx, key, data, ttl)
	return nil
}

// Invalidate deletes keys matching the given patterns.
// Patterns should include "*" glob, e.g. "articles:list:*".
func Invalidate(ctx context.Context, patterns ...string) error {
	if RDB == nil {
		return nil
	}
	for _, pattern := range patterns {
		var cursor uint64
		for {
			keys, nextCursor, err := RDB.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				break
			}
			if len(keys) > 0 {
				RDB.Del(ctx, keys...)
			}
			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}
	return nil
}

// Delete removes specific keys.
func Delete(ctx context.Context, keys ...string) error {
	if RDB == nil || len(keys) == 0 {
		return nil
	}
	return RDB.Del(ctx, keys...).Err()
}

// IncrBy atomically increments a key by delta. Returns the new value.
func IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	if RDB == nil {
		return 0, fmt.Errorf("redis not connected")
	}
	return RDB.IncrBy(ctx, key, delta).Result()
}

// HIncrBy atomically increments a hash field by delta.
func HIncrBy(ctx context.Context, key, field string, delta int64) error {
	if RDB == nil {
		return nil
	}
	return RDB.HIncrBy(ctx, key, field, delta).Err()
}

// HGetAll returns all fields and values in a hash.
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if RDB == nil {
		return nil, fmt.Errorf("redis not connected")
	}
	return RDB.HGetAll(ctx, key).Result()
}
