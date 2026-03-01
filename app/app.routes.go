package app

import (
	"github.com/awesome-goose/goose/modules/router"
)

var (
	ROUTES = router.ForRoutes(
		// Default route shows help
		router.Cli("/", []any{AppController{}, "Help"}),
		// Build command
		router.Cli("/build", []any{AppController{}, "Build"}),
		// Serve command
		router.Cli("/serve", []any{AppController{}, "Serve"}),
		// Help command
		router.Cli("/help", []any{AppController{}, "Help"}),
	)
)
