package modules

import (
	"strconv"
	"strings"
	"time"
)

/*
function to modify the apparence of the locations from the api.
The underscores are replaced with spaces and the cities and countries names are capitalized
*/
func changeLocationsCaracters(locations []string, i int) string {
	location := []rune(locations[i])

	for j := 0; j < len(location); j++ {
		/* The first letter, letters after underscores and the first letter after the dash are capitalized. */
		if (j == 0) || (j > 0 && (location[j-1] == '_') || (location[j-1] == '-')) {
			if (location[j] >= 'a') && (location[j] <= 'z') {
				location[j] = location[j] - 32
			}
		}
	}
	/* The new location is return. strings.ReplaceAll replace all the underscores with spaces using the new string capitalized */
	return strings.ReplaceAll(string(location), "_", " ")
}

/* function that create an array with all the month in letters */
func monthInString() []string {
	var monthString []string

	for i := 1; i <= 12; i++ {
		monthString = append(monthString, time.Month(i).String()) // time.Month convert an int to a month in letters
	}

	return monthString
}

/* function to convert the month string (the numbers) in int */
func stringToInt(stringToReplace string) (int, error) {
	k, err := strconv.Atoi(stringToReplace) // strconv.Atoi convert string to int

	if err != nil {
		return 0, err
	}

	return k, nil
}

/*
function to modify the apparence of the dates from the api.
The month are wrote in letters and the first letter are now in lowercase. The hyphen are now spaces.
*/
func changeDatesCaracters(dates []string, i int) string {
	monthString := monthInString()              // The month array
	date := dates[i]                            // The date were modifiying
	dateInParts := strings.Split(dates[i], "-") // strings.Split allows to split each part of the date using the hyphen
	newDates := ""

	if len(dateInParts) < 3 { // dateInParts should be smaller than 3 because the dates are separated in 3 and the array start at 0
		return date
	}

	monthToReplace := dateInParts[1] // The month

	k, err := stringToInt(monthToReplace)
	if err != nil || k < 1 || k > 12 { // Test to make sure the number asigned to k is a month
		return date
	}

	/* strings.Replace allows to replace in dates[i] the month in number by the month in letters. k-1 is used because the array start at 0.
	And to make sure that the only the month number change the replacement is done only one time. */
	letterMonth := strings.Replace(dates[i], monthToReplace, monthString[k-1], 1)

	/* letterMonth is parcoured to put the month in lowercase */
	capitalizedDates := []rune(letterMonth)
	for j := 0; j < len(capitalizedDates); j++ {
		if (capitalizedDates[j] >= 'A') && (capitalizedDates[j] <= 'Z') {
			capitalizedDates[j] = capitalizedDates[j] + 32 // a is placed 32 number after the A in the ascii table
			newDates += string(capitalizedDates[j])        // the new letter is add in a string
		} else {
			newDates += string(capitalizedDates[j]) // the unchanged letter is also add to the string
		}
	}
	/* The new date is returned. The hyphen is replaced with spaces in the new string (with the month in letter and in lowercase) using strings.ReplaceAll */
	return strings.ReplaceAll(newDates, "-", " ")
}
