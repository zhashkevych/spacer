package spacer

import (
	"context"
	"github.com/pkg/errors"
)

// Dumper is an interface describing DBMS client that creates dump files
type Dumper interface {
	Dump(ctx context.Context, file *TempFile) error
}

// Saver is used to save/retrive dump file from remote object storage
type Saver interface {
	Save(ctx context.Context, file *TempFile) (string, error)
}

// Restorer is an interface describing DBMS client that restores DB from provided dump file
type Restorer interface {
	Restore(ctx context.Context, filename string) error
}

// Saver is used to save/retrive dump file from remote object storage
type DumpDownloader interface {
	GetLatest(ctx context.Context) (string, error)
}

type DumperRestorer interface {
	Dumper
	Restorer
}

type SaverDownloader interface {
	Saver
	DumpDownloader
}

type Spacer struct {
	dumper DumperRestorer
	saver  SaverDownloader
	enc    *Encryptor
	prefix string
}

func NewSpacer(d DumperRestorer, s SaverDownloader, enc *Encryptor, prefix string) *Spacer {
	return &Spacer{dumper: d, saver: s, enc: enc, prefix: prefix}
}

// Export creates dump and saves it using provided Database and Saver objects
func (s *Spacer) Export(ctx context.Context) (string, error) {
	dumpFile, err := NewTempFile(s.prefix, s.enc)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create dump file")
	}

	defer dumpFile.Remove()

	if err := s.dumper.Dump(ctx, dumpFile); err != nil {
		return "", errors.WithMessage(err, "failed to dump db")
	}

	if err := dumpFile.Encrypt(); err != nil {
		return "", errors.WithMessage(err, "failed to encrypt")
	}

	url, err := s.saver.Save(ctx, dumpFile)
	if err != nil {
		return "", errors.WithMessage(err, "failed to save")
	}

	return url, nil
}

// Restore fetches latest dump from object storage using,
func (s *Spacer) Restore(ctx context.Context) error {
	return nil
}

