package fs

import (
	"testing"

	_ "github.com/fedragon/ark/testing"
)

func TestWalk(t *testing.T) {
	cases := []struct {
		name     string
		root     string
		expected int
	}{
		{
			name:     "walk returns all media in a directory",
			root:     "./test/testdata/a",
			expected: 1,
		},
		{
			name:     "walk returns all media in a directory and all its subdirectories",
			root:     "./test/testdata",
			expected: 4,
		},
	}

	for _, c := range cases {
		var count int
		for i := range Walk(c.root, []string{"jpg"}) {
			if i.Err != nil {
				t.Errorf("error: %v", i.Err.Error())
			}

			count++
		}

		if count != c.expected {
			t.Errorf("%v\n\tExpected %v but got %v instead", c.name, c.expected, count)
		}
	}
}
