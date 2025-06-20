package main

import (
	"torrentor/backend/containers/app_container"
	"torrentor/backend/settings"

	"github.com/pkg/errors"

	"github.com/teadove/teasutils/utils/logger_utils"
)

func main() {
	ctx := logger_utils.NewLoggedCtx()

	container, err := app_container.Build(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed to build app container"))
	}

	app := container.WebPresentation.BuildApp()

	err = app.Listen(settings.Settings.WebServer.URL)
	if err != nil {
		panic(errors.Wrap(err, "failed to run app container"))
	}
}
