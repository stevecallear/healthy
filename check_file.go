package healthy

import (
	"context"
	"os"
)

func File(path string) Check {
	return NewCheck(func(ctx context.Context) error {
		_, err := os.Stat(path)
		return err
	}, "type", "file", "target", path)
}
