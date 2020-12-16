package spacer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const filenameTemplate = "%s.dump_%s.sql"

// TempFile is used to create temporary dump files
type TempFile struct {
	encryptor *Encryptor
	file      *os.File
	name      string
	prefix    string
}

func NewTempFile(prefix string, enc *Encryptor) (*TempFile, error) {
	file, err := ioutil.TempFile("", "temp.*.sql")
	if err != nil {
		return nil, err
	}

	return &TempFile{
		encryptor: enc,
		file:      file,
		name:      file.Name(),
		prefix:    prefix,
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

func (f *TempFile) Encrypt() error {
	fileData, err := ioutil.ReadAll(f.file)
	if err != nil {
		return err
	}

	encrypted, err := f.encryptor.Encrypt(fileData)
	if err != nil {
		return err
	}

	if err := f.Remove(); err != nil {
		return err
	}

	filename := f.generateFilename()
	f.file, err = os.Create(filename)
	if err != nil {
		return err
	}
	f.name = filename

	return ioutil.WriteFile(filename, encrypted, 0777)
}

func (f *TempFile) Remove() error {
	return os.Remove(f.name)
}

func (f *TempFile) generateFilename() string {
	return fmt.Sprintf(filenameTemplate, f.prefix, time.Now().Format(time.RFC3339))
}