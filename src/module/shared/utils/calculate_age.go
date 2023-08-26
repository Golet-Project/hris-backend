package utils

import "time"

func CalculateAge(birthdate time.Time) int {
	currentDate := time.Now()
	age := currentDate.Year() - birthdate.Year()

	// Adjust age if the birthdate hasn't occurred yet this year
	if currentDate.YearDay() < birthdate.YearDay() {
		age--
	}

	return age
}
