package main

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/app"
)

func main() {
	setup := config.MustLoad()

	app.Run(setup)
}
