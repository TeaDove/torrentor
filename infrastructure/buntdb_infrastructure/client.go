package buntdb_infrastructure

import (
	"io/fs"
	"os"
	"path"
	"torrentor/settings"

	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

func NewClientFromSettings() (*buntdb.DB, error) {
	err := os.MkdirAll(path.Dir(settings.Settings.BuntDB.DataFile), fs.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "error creating buntdb data directory")
	}

	db, err := buntdb.Open(settings.Settings.BuntDB.DataFile)
	if err != nil {
		return nil, errors.Wrap(err, "error opening buntdb")
	}

	return db, nil
}
