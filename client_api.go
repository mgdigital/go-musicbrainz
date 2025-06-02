package musicbrainz

import (
	"context"
	"strings"
	"time"
)

type clientAPI struct {
	requester requester
}

func (c clientAPI) Artist(ctx context.Context, id string) (Record[Artist], error) {
	return clientAPIByID[Artist](
		ctx,
		c.requester,
		"artist",
		id,
		[]string{"tags", "url-rels"},
	)
}

func (c clientAPI) Release(ctx context.Context, id string) (Record[Release], error) {
	return clientAPIByID[Release](
		ctx,
		c.requester,
		"release",
		id,
		[]string{"artists", "labels", "recordings", "release-groups", "url-rels"},
	)
}

func (c clientAPI) ReleaseGroup(ctx context.Context, id string) (Record[ReleaseGroup], error) {
	return clientAPIByID[ReleaseGroup](
		ctx,
		c.requester,
		"release-group",
		id,
		[]string{"artists", "genres", "url-rels"},
	)
}

func clientAPIByID[T any](ctx context.Context, requester requester, entity, id string, inc []string) (Record[T], error) {
	var data T
	_, err := requester.request(ctx, "/"+entity+"/"+id, map[string]string{
		"inc": strings.Join(inc, "+"),
		"fmt": "json",
	}, &data)
	return Record[T]{
		Date: time.Now(),
		Data: data,
	}, err
}

func (c clientAPI) SearchRelease(ctx context.Context, req SearchReleaseRequest) (SearchReleaseResult, error) {
	var result SearchReleaseResult
	_, err := c.requester.request(ctx, "/release", map[string]string{
		"query": req.Query(),
		"fmt":   "json",
	}, &result)
	return result, err
}

func (c clientAPI) SearchReleaseGroup(ctx context.Context, req SearchReleaseGroupRequest) (SearchReleaseGroupResult, error) {
	var result SearchReleaseGroupResult
	_, err := c.requester.request(ctx, "/release-group", map[string]string{
		"query": req.Query(),
		"fmt":   "json",
	}, &result)
	return result, err
}
