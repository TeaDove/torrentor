package sqlite_infrastructure

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"path"
	"torrentor/settings"

	"github.com/pkg/errors"
)

func NewClientFromSettings() (*gorm.DB, error) {
	err := os.MkdirAll(path.Dir(settings.Settings.SQLite.DataFile), fs.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "error creating buntdb data directory")
	}

	db, err := gorm.Open(sqlite.Open(settings.Settings.SQLite.DataFile), &gorm.Config{})

	return db, nil
}
