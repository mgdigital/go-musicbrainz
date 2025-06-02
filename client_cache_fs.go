package musicbrainz

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
	"time"
)

type clientCacheFS struct {
	Client
	baseDir string
}

func (c *clientCacheFS) Release(ctx context.Context, id string) (Record[Release], error) {
	return getClientCacheFSEntity[Release](ctx, c.baseDir, "release", id, c.Client.Release)
}

func (c *clientCacheFS) Artist(ctx context.Context, id string) (Record[Artist], error) {
	return getClientCacheFSEntity[Artist](ctx, c.baseDir, "artist", id, c.Client.Artist)
}

func (c *clientCacheFS) ReleaseGroup(ctx context.Context, id string) (Record[ReleaseGroup], error) {
	return getClientCacheFSEntity[ReleaseGroup](ctx, c.baseDir, "releasegroup", id, c.Client.ReleaseGroup)
}

func getClientCacheFSEntity[T any](
	ctx context.Context,
	baseDir, name, id string,
	fn func(ctx context.Context, id string) (Record[T], error),
) (Record[T], error) {
	var (
		bytes  []byte
		record Record[T]
		err    error
	)

	filePath := path.Join(baseDir, name, id+".json")
	bytes, err = os.ReadFile(filePath)

	if errors.Is(err, os.ErrNotExist) {
		if record, err = fn(ctx, id); err == nil {
			record.Date = time.Now()
			if bytes, err = json.MarshalIndent(record, "", "  "); err == nil {
				if err = os.MkdirAll(path.Dir(filePath), fs.ModeDir|fs.ModePerm); err == nil {
					err = os.WriteFile(filePath, bytes, 0777)
				}
			}
		}
	} else if err == nil {
		err = json.Unmarshal(bytes, &record)
	}

	return record, err
}
