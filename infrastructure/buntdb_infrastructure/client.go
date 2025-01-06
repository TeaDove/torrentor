package buntdb_infrastructure

import (
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
	"torrentor/settings"
)

func NewClientFromSettings() (*buntdb.DB, error) {
	db, err := buntdb.Open(settings.Settings.BuntDB.DataFile)
	if err != nil {
		return nil, errors.Wrap(err, "error opening buntdb")
	}

	return db, nil
}
