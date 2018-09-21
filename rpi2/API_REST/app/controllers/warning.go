package controllers

import (
	"log"
	"os/exec"

	"app/models"
)

// FireAlarm starts a loud and visual alarm
func FireAlarm(cpd *models.CPD, alarmSoundPath string) {
	go BlinkingAlarm(cpd)
	go soundAlarm(cpd, alarmSoundPath)
}

// soundAlarm plays an alarm sound six times (~15 seconds) of until
// there is not warning (what happens before)
func soundAlarm(cpd *models.CPD, alarmSoundPath string) {
	cmdName := "omxplayer"
	cmdArgs := []string{alarmSoundPath}
	// TODO loop
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		log.Printf("ERROR warning/soundAlarm(): %v\n", err.Error())
	} else {
		log.Printf("Playing sound warning/soundAlarm() %v\n", cmdOut)
	}
}
