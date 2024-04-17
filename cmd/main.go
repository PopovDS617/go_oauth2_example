package main

import (
	"golangoauth2example/internal/app"
	"log"
)

func main() {

	app := app.NewApp()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
