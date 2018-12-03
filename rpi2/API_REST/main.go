package main

import (
	"rpi2/API_REST/app"
)

func main() {
	app := app.App{}
	app.Initialize()
	app.Run()
}
