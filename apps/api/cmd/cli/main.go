package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()

	app := &cli.App{
		Name:  "inari",
		Usage: "photo organiser",
		Action: func(*cli.Context) error {
			logger.Infow("hello!",
				"context", 123,
			)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Errorw("failed to run cli app",
			"err", err)
	}

}
