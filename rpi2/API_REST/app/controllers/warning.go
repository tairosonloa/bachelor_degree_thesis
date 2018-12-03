package controllers

import (
	"log"
	"os/exec"

	"rpi2/API_REST/app/models"
)

// FireAlarm starts a loud and visual alarm
func FireAlarm(cpd *models.CPD, alarmSoundPath string) {
	go BlinkingAlarm(cpd)
	go soundAlarm(cpd, alarmSoundPath)
}

// soundAlarm plays an alarm sound 4 times (~8 seconds) of until
// there is not warning (what happens before)
func soundAlarm(cpd *models.CPD, alarmSoundPath string) {
	cmdName := "omxplayer"
	cmdArgs := []string{alarmSoundPath}
	var cmdOut []byte
	var err error
	for i := 0; cpd.IsWarning() && i < 4; i++ {
		cmdOut, err = exec.Command(cmdName, cmdArgs...).Output()
		if err != nil {
			log.Printf("ERROR warning/soundAlarm(): %v\n%v\n", string(cmdOut), err.Error())
			return
		}
	}
}
