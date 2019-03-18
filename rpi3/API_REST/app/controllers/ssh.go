package controllers

import (
	"log"

	"github.com/sfreiberg/simplessh"
)

// AskOccupation asks control server for classroom occupation. Returns
// a string with the output. Returns an empty string in case of error
func AskOccupation(server, command string) *[]byte {
	client, err := simplessh.ConnectWithKeyFile(server, "root", "")
	if err != nil {
		log.Printf("ERROR controllers.ssh/AskOccupation(): %v\n", err.Error())
		return nil
	}
	defer client.Close()
	output, err := client.Exec(command)
	if err != nil {
		log.Printf("ERROR controllers.ssh/AskOccupation(): %v\n", err.Error())
		return nil
	}
	return &output
}
