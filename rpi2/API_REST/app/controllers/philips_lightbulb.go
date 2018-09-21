package controllers

import (
	"log"
	"time"

	"app/models"

	"github.com/collinux/GoHue"
)

var (
	// bridge represents the hue bridge
	bridge hue.Bridge
	// light represents the lightbulb we have installed
	light hue.Light
)

// InitializeHue connects to hue bridge
func InitializeHue(hueBridgeAddress string, hueBridgeToken string) {
	log.Println("Connecting to Philips Hue Bridge")
	const lightID = 5
	bridge = hue.Bridge{IPAddress: hueBridgeAddress}
	bridge.Login(hueBridgeToken)
	var err error
	light, err = bridge.GetLightByIndex(lightID)
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error())
	} else {
		log.Println("Connected to Philips Hue Bridge")
	}
}

// LightON turns on the philips hue wifi lightbulb on a default white state
func LightON() {
	light.On()
	light.SetBrightness(25)
	light.SetColor(hue.WHITE)
}

// LightOff turns off the philips hue wifi lightbulb
func LightOff() {
	light.Off()
}

// IsLightOn returns true if light inside CPD is on, returns false otherwise
func IsLightOn() bool {
	return light.State.On
}

// BlinkingAlarm set the lightbulb on red blinking to inform of in an alarm status
func BlinkingAlarm(cpd *models.CPD) {
	light.SetColor(hue.RED)
	const blinkMax = 100 // Percent brightness
	const blinkMin = 0   // Percent brightness
	const seconds = 3    // Seconds per blick cycle

	for cpd.IsWarning() {
		light.SetBrightness(blinkMax)
		time.Sleep(time.Second)
		LightOff()
		time.Sleep(time.Second)
	}
}
