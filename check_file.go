package healthy

import (
	"context"
	"os"
)

func File(path string) MetadataCheck {
	return WithMetadata(func(ctx context.Context) error {
		_, err := os.Stat(path)
		return err
	}, "type", "file", "target", path)
}
