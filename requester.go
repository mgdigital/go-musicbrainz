package musicbrainz

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type requester interface {
	request(
		ctx context.Context,
		url string,
		queryParams map[string]string,
		result interface{},
	) (*resty.Response, error)
}
