package spacer

import (
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
)

// Dumper is an interface describing DBMS client that creates dump files
type Dumper interface {
	Dump(ctx context.Context, filename string) error
}

// Saver is used to save/retrive dump file from remote object storage
type Saver interface {
	Save(ctx context.Context, file *DumpFile, folder string) (string, error)
}

// Restorer is an interface describing DBMS client that restores DB from provided dump file
type Restorer interface {
	Restore(ctx context.Context, filename string) error
}

// Saver is used to save/retrive dump file from remote object storage
type Downloader interface {
	GetLatest(ctx context.Context, prefix, folder string) (*DumpFile, error)
}

// DumpRestorer is the interface that groups basic Dumper and Restorer interfaces
type DumpRestorer interface {
	Dumper
	Restorer
}

// SaveDownloader is the interface that groups basic Saver and Downloader interfaces
type SaveDownloader interface {
	Saver
	Downloader
}

type Spacer struct {
	dumper DumpRestorer
	saver  SaveDownloader
	enc    *Encryptor
}

func NewSpacer(d DumpRestorer, s SaveDownloader, enc *Encryptor) *Spacer {
	return &Spacer{dumper: d, saver: s, enc: enc}
}

// Export creates dump and saves it using provided Database and Saver objects
func (s *Spacer) Export(ctx context.Context, prefix, folder string) (string, error) {
	dumpFile, err := NewDumpFile(prefix)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create dump file")
	}

	defer dumpFile.Remove()

	if err := s.dumper.Dump(ctx, dumpFile.Name()); err != nil {
		return "", errors.WithMessage(err, "failed to dump db")
	}

	if err := s.encryptFile(dumpFile); err != nil {
		return "", errors.WithMessage(err, "failed to encrypt file")
	}

	url, err := s.saver.Save(ctx, dumpFile, folder)
	if err != nil {
		return "", errors.WithMessage(err, "failed to save")
	}

	return url, nil
}

// Restore fetches latest dump from object storage using,
func (s *Spacer) Restore(ctx context.Context, prefix, folder string) error {
	file, err := s.saver.GetLatest(ctx, prefix, folder)
	if err != nil {
		return errors.WithMessage(err, "failed to get latest dump file")
	}

	defer file.Remove()

	if err := s.decryptFile(file); err != nil {
		return errors.WithMessage(err, "failed to decrypt file")
	}

	return s.dumper.Restore(ctx, file.Name())
}

func (s *Spacer) encryptFile(f *DumpFile) error {
	fileData, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return err
	}

	encrypted, err := s.enc.Encrypt(fileData)
	if err != nil {
		return err
	}

	return f.Write(encrypted)
}

func (s *Spacer) decryptFile(f *DumpFile) error {
	fileData, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return err
	}

	decrypted, err := s.enc.Decrypt(fileData)
	if err != nil {
		return err
	}

	return f.Write(decrypted)
}