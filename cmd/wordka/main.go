package main

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/app"
)

func main() {
	env, cfg := config.MustLoad()

	app.Run(env, cfg)
}
