package healthy_test

import (
	"context"
	"maps"
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	"github.com/stevecallear/healthy"
)

func TestFile_Check(t *testing.T) {
	t.Run("should return an error if the file does not exist", func(t *testing.T) {
		fn, close := tempFile()
		close() // remove immediately

		sut := healthy.File(fn)
		err := sut.Healthy(context.Background())
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should return nil if the file exists", func(t *testing.T) {
		fn, close := tempFile()
		defer close()

		f, err := os.Create(fn)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		defer func() {
			os.Remove(fn)
		}()

		sut := healthy.File(fn)
		err = sut.Healthy(context.Background())
		if err != nil {
			t.Errorf("got %v, expected nil", err)
		}
	})
}

func TestFile_Info(t *testing.T) {
	t.Run("should return the check info", func(t *testing.T) {
		const file = "test.txt"
		exp := healthy.Metadata{"type": "file", "target": file}
		act := healthy.File(file).Metadata()
		if !maps.Equal(act, exp) {
			t.Errorf("got %v, expected %v", act, exp)
		}
	})
}

func tempFile() (string, func()) {
	fn := randString(16) + ".tmp"
	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	return fn, func() {
		defer f.Close()
		defer os.Remove(fn)
	}
}

func randString(l int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := new(strings.Builder)
	for range l {
		b.WriteByte(alphabet[rand.Int()%len(alphabet)])
	}
	return b.String()
}
