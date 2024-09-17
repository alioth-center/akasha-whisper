package global

import (
	"github.com/alioth-center/akasha-whisper/app/dao"
	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

var (
	LoginCookieCacheInstance         cache.Cache
	BearerTokenBloomFilterInstance   *dao.BearerTokenBloomFilter
	OpenaiClientCacheInstance        concurrency.Map[int, openai.Client]
	OpenaiClientSecretsCacheInstance concurrency.Map[int, *openai.Config]
)
