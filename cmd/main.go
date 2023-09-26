package main

import (
	"godev/internal/app"
	"godev/pkg/embed"
)

func main() {
	a := app.NewApp()
	embed.InitEmbed()
	a.Start()
}
