package main

import (
	"rpi3/API_REST/app"
)

func main() {
	app := app.App{}
	app.Initialize()
	app.Run()
}
