package controllers

import (
	"log"

	"github.com/collinux/GoHue"
)

const (
	lightID = 5
)

var (
	bridge hue.Bridge
	light  hue.Light
)

// InitializeHue connects to hue bridge
func InitializeHue(hueBridgeAddress string, hueBridgeToken string) {
	log.Println("Connecting to Philips Hue Bridge")
	bridge = hue.Bridge{IPAddress: hueBridgeAddress}
	bridge.Login(hueBridgeToken)
	var err error
	light, err = bridge.GetLightByIndex(lightID)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	} else {
		log.Println("Connected to Philips Hue Bridge")
	}
}

// LightON turn on the philips hue wifi lightbulb on a default white state
func LightON() {
	light.On()
	light.SetBrightness(25)
	light.SetColor(hue.WHITE)
}

// LightOff turn off the philips hue wifi lightbulb
func LightOff() {
	light.Off()
}
