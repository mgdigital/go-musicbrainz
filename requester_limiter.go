package musicbrainz

import (
	"context"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

type requesterLimiter struct {
	requester requester
	limiter   *rate.Limiter
}

func (r requesterLimiter) request(
	ctx context.Context,
	url string,
	queryParams map[string]string,
	result interface{},
) (*resty.Response, error) {
	if err := r.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return r.requester.request(ctx, url, queryParams, result)
}
