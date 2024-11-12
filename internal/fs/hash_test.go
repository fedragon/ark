package fs

import (
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/fedragon/ark/testing"
)

func TestHash(t *testing.T) {
	dest := "./test/testdata"

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
		a, err := Hash(c.pathA)
		if err != nil {
			t.Errorf(err.Error())
		}
		b, err := Hash(c.pathB)
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
	for i := 0; i < b.N; i++ {
		_, err := Hash(filepath.Join("./test/data", "doge.jpg"))
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}
