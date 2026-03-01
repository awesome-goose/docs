package main

import (
	"os"

	"github.com/awesome-goose/docs/app"
	"github.com/awesome-goose/goose"
	"github.com/awesome-goose/goose/log"
	"github.com/awesome-goose/goose/log/formatters"
	"github.com/awesome-goose/goose/log/modifiers"
	"github.com/awesome-goose/goose/log/processors"
	"github.com/awesome-goose/goose/platforms/cli"
	"github.com/awesome-goose/goose/types"
)

func main() {
	stop, err := goose.Start(
		goose.CLI(cliPlatform, rootModule, initializers),
	)
	if err != nil {
		panic(err)
	}

	defer stop()
}

var (
	cliPlatform = cli.NewPlatform(
		cli.WithName("goose-docs"),
	)
	rootModule   = &app.AppModule{}
	initializers = []func(container types.Container) error{
		func(container types.Container) error {
			return container.Register(
				func() types.Log {
					return log.NewLog(
						log.AppLogChannel("std"),
						log.NewLogger(
							[]types.Modifier{
								modifiers.NewUUID(),
								modifiers.NewColorTagsModifier(),
								modifiers.NewSystemInfo(),
								modifiers.NewStackTrace(),
							},
							formatters.NewSyslog("goose-docs", os.Getpid()),
							processors.NewConsole(),
						),
					)
				},
				"",
				true,
			)
		},
	}
)
