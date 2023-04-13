package fs

import (
	"testing"

	. "github.com/fedragon/ark/testing"

	"github.com/spf13/afero"
)

func TestWalk(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := CopyTestData(fs, "/src/test/data", "doge.jpg", "same-doge.jpg", "grumpy-cat.jpg"); err != nil {
		t.Fatalf(err.Error())
	}

	cases := []struct {
		name     string
		root     string
		expected int
	}{
		{
			name:     "walk returns all media in a directory",
			root:     "/src/test/data",
			expected: 3,
		},
		{
			name:     "walk returns all media in a directory and all its subdirectories",
			root:     "/src",
			expected: 3,
		},
	}

	for _, c := range cases {
		var count int
		for i := range Walk(fs, c.root, []string{".jpg"}) {
			if i.Err != nil {
				t.Errorf(i.Err.Error())
			}

			count++
		}

		if count != c.expected {
			t.Errorf("%v\n\tExpected %v but got %v instead", c.name, c.expected, count)
		}
	}
}
