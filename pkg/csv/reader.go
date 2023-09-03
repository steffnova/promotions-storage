package csv

import (
	"io"
	"os"
)

// Reader returns new [io.Reader] or error
type Reader func() (io.Reader, error)

// FileReader returns file reader. The fileName parameter
// defines name of a file that will be used by reader.
func FileReader(fileName string) Reader {
	return func() (io.Reader, error) {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
