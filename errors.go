package musicbrainz

import (
	"errors"
	"fmt"
)

var (
	Err         = errors.New("musicbrainz client")
	ErrNotFound = fmt.Errorf("%w: not found", Err)
)
