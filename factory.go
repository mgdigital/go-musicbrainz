package musicbrainz

import (
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/navidrome/navidrome/conf"
	"golang.org/x/time/rate"
)

func NewClient(config Config) Client {
	var client Client

	client = &clientAPI{
		requester: requesterLimiter{
			requester: requesterResty{
				resty: resty.New().
					SetBaseURL(config.BaseURL).
					SetTimeout(config.Timeout).
					SetHeader("User-Agent", config.UserAgent).
					SetRetryCount(config.RetryCount).
					SetRetryWaitTime(config.RetryWaitTime).
					SetRetryMaxWaitTime(config.RetryMaxWaitTime),
			},
			limiter: rate.NewLimiter(config.RateLimit, config.RateBurst),
		},
	}

	if config.FSCacheConfig.BaseDir != "" {
		client = &clientCacheFS{
			Client:  client,
			baseDir: filepath.Join(conf.Server.DataFolder, "musicbrainz"),
		}
	}

	if config.LRUCacheConfig.Size > 0 {
		client = &clientCacheInMem{
			Client: client,
			lru: expirable.NewLRU[string, any](
				config.LRUCacheConfig.Size,
				nil,
				config.LRUCacheConfig.TTL,
			),
		}
	}

	return client
}
