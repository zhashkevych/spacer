package spacer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const filenameTemplate = "dump_%d.sql"

// DumpFile is used to create temporary dump files
type DumpFile struct {
	file *os.File
	name string
}

func NewDumpFile(prefix string) (*DumpFile, error) {
	filename := generateFilename(prefix)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &DumpFile{
		file: file,
		name: filename,
	}, nil
}

func (f *DumpFile) Name() string {
	return f.name
}

func (f *DumpFile) Size() (int64, error) {
	stat, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (f *DumpFile) Reader() io.Reader {
	return f.file
}

func (f *DumpFile) Write(data []byte) error {
	return ioutil.WriteFile(f.name, data, 0777)
}

func (f *DumpFile) Remove() error {
	return os.Remove(f.name)
}

func generateFilename(prefix string) string {
	filename := fmt.Sprintf(filenameTemplate, time.Now().Unix())
	if prefix == "" {
		return filename
	}

	return fmt.Sprintf("%s.%s", prefix, filename)
}
