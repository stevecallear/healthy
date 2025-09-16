package healthy

import (
	"context"
	"os"
)

// File returns a file health check
// The check returns nil if the file exists.
func File(path string) MetadataCheck {
	return WithMetadata(func(ctx context.Context) error {
		_, err := os.Stat(path)
		return err
	}, "type", "file", "target", path)
}
