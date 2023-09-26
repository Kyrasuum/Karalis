package main

import (
	"karalis/internal/app"
	"karalis/pkg/embed"
)

func main() {
	a := app.NewApp()
	embed.InitEmbed()
	a.Start()
}
