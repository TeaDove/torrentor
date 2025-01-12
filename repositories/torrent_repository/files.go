package torrent_repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"path"
	"strings"
)

func makeFileToPathKey(id uuid.UUID) string {
	return fmt.Sprintf("file:%s", id)
}

func makeFileData(torrentInfoHash string, fileName string) string {
	return fmt.Sprintf("%s::%s", torrentInfoHash, fileName)
}

func parseFileData(v string) (string, string, error) {
	fields := strings.Split(v, "::")
	if len(fields) != 2 {
		return "", "", errors.New("invalid file data")
	}

	return fields[0], fields[1], nil
}

func (r *Repository) saveFiles(torrent *Torrent) error {
	files := torrent.FlatFiles()
	for _, file := range files {
		_, _, err := r.db.Set(makeFileToPathKey(file.Id), makeFileData(torrent.InfoHash, file.Path), nil)
		if err != nil {
			return errors.Wrap(err, "failed to save torrent")
		}
	}

	return nil
}

func (r *Repository) FileGetPath(_ context.Context, id uuid.UUID) (string, error) {
	val, err := r.db.Get(makeFileToPathKey(id))
	if err != nil {
		return "", errors.Wrap(err, "failed to get file by path")
	}

	hash, name, err := parseFileData(val)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse file data")
	}

	return path.Join(r.torrentDataDir, hash, name), nil
}
