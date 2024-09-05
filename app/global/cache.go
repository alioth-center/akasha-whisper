package global

import (
	"github.com/alioth-center/akasha-whisper/app/dao"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

var (
	BearerTokenBloomFilterInstance *dao.BearerTokenBloomFilter
	OpenaiClientCacheInstance      concurrency.Map[int, openai.Client]
)
