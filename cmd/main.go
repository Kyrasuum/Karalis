package main

import (
	"fmt"

	"karalis/internal/app"
)

func main() {
	a := app.NewApp()
	if a != nil {
		err := a.Start()
		if err != nil {
			fmt.Printf("ERR: %+v\n", err)
		}
	}
}
