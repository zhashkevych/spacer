package spacer

import "context"

// Restorer is an interface describing DBMS client that restores DB from provided dump file
type Restorer interface {
	Restore(filename string) error
}

// Saver is used to save/retrive dump file from remote object storage
type DumpDownloader interface {
	GetLatest(ctx context.Context) (string, error)
}

// Restore fetches latest dump from object storage using,
func Restore(r Restorer, dd DumpDownloader, enc *Encryptor) (string, error) {
	return "", nil
}

