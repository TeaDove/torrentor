package views

import (
	"embed"

	"github.com/pkg/errors"
)

//go:embed *
var Static embed.FS

func init() {
	files, err := Static.ReadDir(".")
	if err != nil {
		panic(errors.Wrap(err, "failed to read static directory"))
	}

	if len(files) == 0 {
		panic(errors.New("no static files found"))
	}
}
