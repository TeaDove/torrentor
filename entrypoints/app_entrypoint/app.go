package main

import (
	"torrentor/containers/tg_bot_container"

	"github.com/teadove/teasutils/utils/di_utils"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func main() {
	ctx := logger_utils.NewLoggedCtx()

	container, err := di_utils.BuildFromSettings(ctx, tg_bot_container.Build)
	if err != nil {
		panic(err)
	}

	container.TGBotPresentation.PollerRun(ctx)
}

// http.Handle("/", http.FileServer(http.Dir("/Users/teadove/Downloads/Shameless.S02.720p.BDRip.x264.ac3.rus.eng")))
//
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		fmt.Println("Error starting server:", err)
//		os.Exit(1)
//	}
