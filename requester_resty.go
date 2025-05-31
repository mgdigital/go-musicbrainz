package musicbrainz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type requesterResty struct {
	resty *resty.Client
}

func (r requesterResty) request(
	ctx context.Context,
	url string,
	queryParams map[string]string,
	result interface{},
) (*resty.Response, error) {
	res, err := r.resty.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&result).
		Execute(http.MethodGet, url)
	if err == nil {
		if !res.IsSuccess() {
			if res.StatusCode() == 404 {
				err = ErrNotFound
			} else {
				err = fmt.Errorf("%w: %s", Err, res.Status())
			}
		}
	}
	return res, err
}
