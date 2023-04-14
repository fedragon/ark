package testing

import (
	"os"
	"path"
	"runtime"
)

// Change the working directory to the project root: useful when reading files for testing purposes.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
