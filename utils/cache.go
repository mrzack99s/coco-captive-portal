package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9"
)

var redisCache *redis.Client

func RedisCache() *redis.Client {
	return redisCache
}

func SetupCache() {
	redisCache = redis.NewClient(&redis.Options{
		Addr:               "127.0.0.1:6379",
		Password:           "",
		DB:                 0,
		IdleCheckFrequency: time.Second * 29,
	})
}

func CacheGetWithRawKey(key string, v interface{}) (err error) {
	var value string
	value, err = redisCache.Get(context.Background(), key).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), v)
	return
}

func CacheGet(prefix, key string, v interface{}) (err error) {
	var value string
	value, err = redisCache.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), v)
	return
}

func CachePop(stack string, v interface{}) (err error) {
	var value string
	value, err = redisCache.LPop(context.Background(), stack).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), v)
	return
}

func CachePush(stack string, v interface{}) (err error) {
	var value []byte
	if reflect.TypeOf(v).Name() != "string" {
		value, err = json.Marshal(v)
		if err != nil {
			return
		}

		err = redisCache.LPush(context.Background(), stack, string(value)).Err()
		if err != nil {
			return
		}
	} else {
		err = redisCache.LPush(context.Background(), stack, v).Err()
		if err != nil {
			return
		}
	}
	return
}

func CacheGetString(prefix, key string) (value string, err error) {
	value, err = redisCache.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Result()
	return
}

func CacheSet(prefix, key string, v interface{}) (err error) {
	var value []byte

	if reflect.TypeOf(v).Name() != "string" {
		value, err = json.Marshal(v)
		if err != nil {
			return
		}

		err = redisCache.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), string(value), 0).Err()
		if err != nil {
			return
		}
	} else {
		err = redisCache.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), v, 0).Err()
		if err != nil {
			return
		}
	}

	return
}

func CacheSetWithTimeDuration(prefix, key string, v interface{}, t time.Duration) (err error) {
	var value []byte

	if reflect.TypeOf(v).Name() != "string" {
		value, err = json.Marshal(v)
		if err != nil {
			return
		}

		err = redisCache.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), string(value), t).Err()
		if err != nil {
			return
		}
	} else {
		err = redisCache.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), v, t).Err()
		if err != nil {
			return
		}
	}

	return
}

func CacheGetAllKey(prefix string) (keys []string, err error) {
	keys, _, err = redisCache.Scan(context.Background(), 0, fmt.Sprintf("%s:*", prefix), 0).Result()
	if err != nil {
		panic(err)
	}
	return
}

func CacheFindExistingRawKey(key string) bool {
	_, err := redisCache.Get(context.Background(), key).Result()
	if err != redis.Nil {
		return true
	} else {
		return false
	}
}

func CacheFindExistingKey(prefix, key string) bool {
	_, err := redisCache.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Result()
	if err != redis.Nil {
		return true
	} else {
		return false
	}
}
func CacheDeleteWithPrefix(prefix string) error {
	iter := redisCache.Scan(context.Background(), 0, fmt.Sprintf("%s:*", prefix), 0).Iterator()
	for iter.Next(context.Background()) {
		val := iter.Val()
		redisCache.Del(context.Background(), val)
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func CacheCountWithPrefix(prefix string) (int, error) {
	keys, _, err := redisCache.Scan(context.Background(), 0, fmt.Sprintf("%s:*", prefix), 0).Result()
	if err != nil {
		return -1, err
	}
	return len(keys), nil
}

func CacheDeleteWithRawKey(key string) (err error) {
	err = redisCache.Del(context.Background(), key).Err()
	return
}

func CacheDelete(prefix, key string) (err error) {
	err = redisCache.Del(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Err()
	return
}
