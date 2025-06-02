package musicbrainz

import (
	"time"

	"golang.org/x/time/rate"
)

type FSCacheConfig struct {
	BaseDir string
}

type LRUCacheConfig struct {
	Size int
	TTL  time.Duration
}

type Config struct {
	FSCacheConfig
	LRUCacheConfig
	BaseURL          string
	UserAgent        string
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	RateLimit        rate.Limit
	RateBurst        int
	Timeout          time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		LRUCacheConfig: LRUCacheConfig{
			Size: 1000,
			TTL:  time.Minute * 10,
		},
		BaseURL:          "https://musicbrainz.org/ws/2",
		UserAgent:        "go-musicbrainz-clientAPI",
		RetryCount:       20,
		RetryWaitTime:    5 * time.Second,
		RetryMaxWaitTime: 60 * time.Second,
		RateLimit:        rate.Every(time.Second),
		RateBurst:        3,
		Timeout:          30 * time.Second,
	}
}
