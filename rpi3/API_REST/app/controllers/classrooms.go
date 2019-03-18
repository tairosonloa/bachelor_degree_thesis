package controllers

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"rpi3/API_REST/app/models"
)

// GetClassroomsStatus returns a struct whit classroom status
// Valid status are: 0 - free, 1 - occupied, 2 - will be occupied in
// the next 30 minutes, 3 - will be occupied in the next 10 minutes
func GetClassroomsStatus(reservations []*models.Reservation) *models.Classrooms {
	classrooms := models.Classrooms{
		F16: 0,
		F18: 0,
		C05: 0,
		C06: 0,
	}
	// Get current hour and minutes
	t := time.Now()
	ch := t.Hour()
	cm := t.Minute()
	// Check classrooms status
	for _, res := range reservations {
		// Check for classrooms that are currently occupied
		if (res.StartHour < ch || (res.StartHour == ch && res.StartMinute <= cm)) && (res.EndHour > ch || (res.EndHour == ch && res.EndMinute > cm)) {
			switch strings.ToLower(res.Classroom) {
			case "4.0.f16":
				classrooms.F16 = 1
			case "4.0.f18":
				classrooms.F18 = 1
			case "2.2.c05":
				classrooms.C05 = 1
			case "2.2.c06":
				classrooms.C06 = 1
			default:
				return nil
			}
			// Check for classrooms that will be occupied in the next 30 minutes
		} else if (res.StartHour == ch && res.StartMinute <= cm+30) || (res.StartHour == ch+1 && 60-cm+res.StartMinute <= 30) {
			switch strings.ToLower(res.Classroom) {
			case "4.0.f16":
				if classrooms.F16 != 1 {
					if 60-cm+res.StartMinute <= 10 { // Check if classroom will be occupied in the next 10 minutes
						classrooms.F16 = 3
					} else {
						classrooms.F16 = 2
					}
				}
			case "4.0.f18":
				if classrooms.F18 != 1 {
					if 60-cm+res.StartMinute <= 10 {
						classrooms.F18 = 3
					} else {
						classrooms.F18 = 2
					}
				}
			case "2.2.c05":
				if classrooms.C05 != 1 {
					if 60-cm+res.StartMinute <= 10 {
						classrooms.C05 = 3
					} else {
						classrooms.C05 = 2
					}
				}
			case "2.2.c06":
				if classrooms.C06 != 1 {
					if 60-cm+res.StartMinute <= 10 {
						classrooms.C06 = 3
					} else {
						classrooms.C06 = 2
					}
				}
			default:
				return nil
			}
		}
	}
	return &classrooms
}

// GetClassroomsOccupation returns a struct with classrooms occupations statistics
func GetClassroomsOccupation(server, command string) *models.Occupation {
	classrooms := [...]string{"f16", "f18", "c05", "c06"}
	occupation := models.Occupation{}
	for _, c := range classrooms {
		// Ask control server for classroom occupation
		output := AskOccupation(server, command+" "+c)
		if output != nil {
			scanner := bufio.NewScanner(bytes.NewReader(*output))
			stats := models.OccupationStats{}
			for scanner.Scan() {
				// Check stats
				if strings.Contains(strings.ToLower(scanner.Text()), "(apagado)") {
					stats.Shutdown++
				} else if strings.Contains(strings.ToLower(scanner.Text()), "(debian)") {
					stats.Linux++
				} else if strings.Contains(strings.ToLower(scanner.Text()), "(windows)") {
					stats.Windows++
				} else if strings.Contains(strings.ToLower(scanner.Text()), "timeout") {
					stats.TimeOut++
				} else if strings.Contains(strings.ToLower(scanner.Text()), "pid comentario") {
					stats.StudentsLinux++
				} else if strings.Contains(strings.ToLower(scanner.Text()), "id. estado") {
					stats.StudentsWindows++
				}
			}
			// Check classroom
			switch strings.ToLower(c) {
			case "f16":
				occupation.F16 = stats
			case "f18":
				occupation.F18 = stats
			case "c05":
				occupation.C05 = stats
			case "c06":
				occupation.C06 = stats
			}
		} else {
			return nil
		}
	}
	return &occupation
}
