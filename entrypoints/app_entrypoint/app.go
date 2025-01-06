package main

import (
	"torrentor/containers/app_container"

	"github.com/teadove/teasutils/utils/di_utils"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func main() {
	ctx := logger_utils.NewLoggedCtx()

	container, err := di_utils.BuildFromSettings(ctx, app_container.Build)
	if err != nil {
		panic(err)
	}

	go container.TGBotPresentation.PollerRun(ctx)
	err = container.WebPresentation.Run(ctx)
	if err != nil {
		panic(err)
	}
}

// http.Handle("/", http.FileServer(http.Dir("/Users/teadove/Downloads/Shameless.S02.720p.BDRip.x264.ac3.rus.eng")))
//
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		fmt.Println("Error starting server:", err)
//		os.Exit(1)
//	}
