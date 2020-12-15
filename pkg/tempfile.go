package spacer

import (
	"fmt"
	"io"
	"os"
	"time"
)

const filenameTemplate = "dump_%d.tar.gz"

// TempFile is used to create temporary dump files
type TempFile struct {
	file *os.File
	name string
}

func NewTempFile() (*TempFile, error) {
	filename := generateFilename()
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &TempFile{
		file: f,
		name: filename,
	}, nil
}

func (f *TempFile) Name() string {
	return f.name
}

func (f *TempFile) Size() (int64, error) {
	stat, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (f *TempFile) Reader() io.Reader {
	return f.file
}

func (f *TempFile) Remove() error {
	return os.Remove(f.name)
}

func generateFilename() string {
	return fmt.Sprintf(filenameTemplate, time.Now().Unix())
}
