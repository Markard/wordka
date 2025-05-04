package main

import (
	"fmt"
	"github.com/Markard/wordka/internal/config"
)

func main() {
	env, cfg := config.MustLoad()
	fmt.Println(env, cfg)
}
