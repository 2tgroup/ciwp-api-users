// redisCache redis provides a redis interface.
package redisCache

import (
	"fmt"
	"time"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"github.com/go-redis/redis"
)

var (
	rc *redis.Client
)

// cache is an implementation of httpcache.Cache that caches responses in a

func init() {
	opt, err := redis.ParseURL(config.DataConfig.Redis["cache"].Host)
	if err != nil {
		panic(err)
	}
	rc = redis.NewClient(opt)
}

// cacheKey modifies an httpcache key for use in redis. Specifically, it
func cacheKey(prefix, key string) string {
	return fmt.Sprintf("%v:%v", prefix, key)
}

// Get returns the response corresponding to key if present.
func Get(prefix, key string) (resp []byte, ok bool) {
	item, err := rc.Get(cacheKey(prefix, key)).Result()
	if err != nil {
		return nil, false
	}
	return []byte(item), true
}

// Set saves a response to the cache as key.
func Set(prefix, key string, resp []byte, expTime int) bool {
	err := rc.Set(cacheKey(prefix, key), resp, time.Minute*time.Duration(expTime)).Err()
	if err != nil {
		return false
	}
	return true
}

//Del removes the response with key from the cache.
func Del(prefix, key string) {
	rc.Del(cacheKey(prefix, key))
}
