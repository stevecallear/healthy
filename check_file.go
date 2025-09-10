package healthy

import (
	"context"
	"os"
)

func File(path string) InfoCheck {
	return NewCheck(func(ctx context.Context) error {
		_, err := os.Stat(path)
		return err
	}, "type", "file", "target", path)
}
