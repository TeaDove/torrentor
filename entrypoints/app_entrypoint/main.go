package main

import (
	"github.com/pkg/errors"
	"torrentor/containers/app_container"

	"github.com/teadove/teasutils/utils/di_utils"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func main() {
	ctx := logger_utils.NewLoggedCtx()

	container, err := di_utils.BuildFromSettings(ctx, app_container.Build)
	if err != nil {
		panic(errors.Wrap(err, "failed to build app container"))
	}

	go container.TGBotPresentation.PollerRun(ctx)
	err = container.WebPresentation.Run(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed to run app container"))
	}
}
