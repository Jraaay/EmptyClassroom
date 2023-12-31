package cache

import (
	"EmptyClassroom/logs"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	GlobalCache *cache.Cache
)

func InitCache() {
	GlobalCache = cache.New(5*time.Minute, 5*time.Minute)
	err := GlobalCache.LoadFile(".ec.gob")
	if err != nil {
		logs.CtxError(nil, "GlobalCache.LoadFile error: %v", err)
		return
	}
}

func GetCache(key string) (interface{}, bool) {
	return GlobalCache.Get(key)
}

func GetCacheWithExpiration(key string) (interface{}, time.Time, bool) {
	return GlobalCache.GetWithExpiration(key)
}

func SetCache(key string, value interface{}, expiration time.Duration) {
	GlobalCache.Set(key, value, expiration)
	err := GlobalCache.SaveFile(".ec.gob")
	if err != nil {
		logs.CtxError(nil, "GlobalCache.SaveFile error: %v", err)
		return
	}
}
