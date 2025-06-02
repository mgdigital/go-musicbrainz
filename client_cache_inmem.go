package musicbrainz

import (
	"context"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type clientCacheInMem struct {
	Client
	lru *expirable.LRU[string, any]
}

func (c *clientCacheInMem) Artist(ctx context.Context, id string) (Record[Artist], error) {
	return getClientCacheInMemEntity[Artist](ctx, c.lru, c.Client.Artist, "artist", id)
}

func (c *clientCacheInMem) Release(ctx context.Context, id string) (Record[Release], error) {
	return getClientCacheInMemEntity[Release](ctx, c.lru, c.Client.Release, "release", id)
}

func (c *clientCacheInMem) ReleaseGroup(ctx context.Context, id string) (Record[ReleaseGroup], error) {
	return getClientCacheInMemEntity[ReleaseGroup](ctx, c.lru, c.Client.ReleaseGroup, "releasegroup", id)
}

func getClientCacheInMemEntity[T any](
	ctx context.Context,
	lru *expirable.LRU[string, any],
	fn func(ctx context.Context, id string) (Record[T], error),
	entity, id string,
) (result Record[T], err error) {

	key := entity + "_" + id
	if rawResult, ok := lru.Get(key); ok {
		return rawResult.(Record[T]), nil
	}
	result, err = fn(ctx, id)
	if err != nil {
		return result, err
	}
	lru.Add(key, result)
	return result, nil
}
