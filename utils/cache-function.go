package utils

import (
	"fmt"

	"github.com/allegro/bigcache/v3"
)

func GetCacheValAsInt64(key string, c *bigcache.BigCache) (v int64, err error) {

	bVal, err := c.Get(key)
	if err != nil {
		return
	}

	v, err = StringToInt64(string(bVal))

	return
}

func GetCacheValAsString(key string, c *bigcache.BigCache) (v string, err error) {

	bVal, err := c.Get(key)
	if err != nil {
		return
	}

	v = string(bVal)

	return
}

func SetCacheFromInt64(key string, v int64, c *bigcache.BigCache) (err error) {

	err = c.Set(key, []byte(fmt.Sprintf("%d", v)))

	return
}
