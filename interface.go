package musicbrainz

import (
	"context"
)

type Client interface {
	Artist(ctx context.Context, id string) (Record[Artist], error)
	Release(ctx context.Context, id string) (Record[Release], error)
	ReleaseGroup(ctx context.Context, id string) (Record[ReleaseGroup], error)
}
