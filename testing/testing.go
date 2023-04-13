package testing

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
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

func CopyTestData(fs afero.Fs, dest string, files ...string) error {
	fs.MkdirAll(dest, 0755)

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, name := range files {
		f, err := os.Open(filepath.Join(workdir, "/test/data/", name))
		if err != nil {
			return err
		}

		if err := afero.WriteReader(fs, filepath.Join(dest, name), bufio.NewReader(f)); err != nil {
			return err
		}
		f.Close()
	}

	return nil
}
