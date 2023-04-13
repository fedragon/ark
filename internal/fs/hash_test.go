package fs

import (
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/fedragon/ark/testing"

	"github.com/spf13/afero"
)

func TestHash(t *testing.T) {
	fs := afero.NewMemMapFs()

	dest := "/src"
	if err := CopyTestData(fs, dest, "doge.jpg", "same-doge.jpg", "grumpy-cat.jpg"); err != nil {
		t.Fatalf(err.Error())
	}

	cases := []struct {
		name     string
		pathA    string
		pathB    string
		expected bool
	}{
		{
			name:     "hashing the same file twice returns the same value",
			pathA:    dest + "/doge.jpg",
			pathB:    dest + "/doge.jpg",
			expected: true,
		},
		{
			name:     "hashing two files with same content but different name returns the same value",
			pathA:    dest + "/doge.jpg",
			pathB:    dest + "/same-doge.jpg",
			expected: true,
		},
		{
			name:     "hashing two different files returns different values",
			pathA:    dest + "/doge.jpg",
			pathB:    dest + "/grumpy-cat.jpg",
			expected: false,
		},
	}

	for _, c := range cases {
		a, err := hash(fs, c.pathA)
		if err != nil {
			t.Errorf(err.Error())
		}
		b, err := hash(fs, c.pathB)
		if err != nil {
			t.Errorf(err.Error())
		}

		equal := reflect.DeepEqual(a, b)
		if equal != c.expected {
			t.Errorf("%v\n\tExpected %v but got %v instead", c.name, c.expected, equal)
		}
	}
}

func BenchmarkHash(b *testing.B) {
	fs := afero.NewMemMapFs()

	dest := "/src"
	name := "doge.jpg"
	if err := CopyTestData(fs, dest, name); err != nil {
		b.Fatalf(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err := hash(fs, filepath.Join(dest, name))
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}
