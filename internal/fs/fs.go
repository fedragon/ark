package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fedragon/ark/internal/db"

	"github.com/djherbis/times"
	"lukechampine.com/blake3"
)

func hash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := blake3.New(256, nil)
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Walk traverses the directory tree rooted at root, sending all media files (with extensions in fileTypes) to the returned channel.
// It spawns a goroutine to walk the tree and immediately returns a read-only channel to receive the values.
// In case of errors, the channel will receive a Media with the Err field set.
func Walk(root string, fileTypes []string) <-chan db.Media {
	media := make(chan db.Media)

	go func() {
		defer close(media)

		typesMap := make(map[string]struct{})
		for _, t := range fileTypes {
			typesMap[t] = struct{}{}
		}

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				t, err := times.Stat(path)
				if err != nil {
					return err
				}

				ext := strings.ToLower(filepath.Ext(d.Name()))
				if _, exists := typesMap[ext]; exists {
					bytes, err := hash(path)
					if err != nil {
						return err
					}

					media <- db.Media{
						Path:    path,
						Hash:    bytes,
						ModTime: t.ModTime(),
					}
				}
			}

			return nil
		})
		if err != nil {
			media <- db.Media{Err: err}
		}
	}()

	return media
}
