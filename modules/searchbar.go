package modules

import (
	"strconv"
	"strings"
)

// SearchArtist - Search for an artist based on a keyword
func (artistsData *General) SearchArtist(search string) ([]Artist, error) {
	var res []Artist

	// Search by creation date
	a, atoiErr := strconv.Atoi(search)
	if atoiErr == nil { // Successful search
		for _, art := range artistsData.Artists {
			if art.CreationDate == a { // Iterate through artists and compare creation dates
				res = append(res, art)
			}
		}
	}

	// Search by artist name
	for _, art := range artistsData.Artists {
		if strings.Contains(strings.ToLower(art.Name), strings.ToLower(search)) {
			res = append(res, art)
			continue
		}

		// Search by first album
		if strings.Contains(strings.ToLower(art.FirstAlbum), strings.ToLower(search)) {
			res = append(res, art)
			continue
		}

		// Search by members
		for _, member := range art.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(search)) {
				res = append(res, art)
				break
			}
		}

		// Search by concert dates in Relations (DatesLocations)
		for _, dates := range art.DatesLocations {
			for _, date := range dates {
				if strings.Contains(date, search) { // Check if the entered date matches
					res = append(res, art)
					break
				}
			}
		}

		// Search by city in Relations (DatesLocations)
		for city := range art.DatesLocations {
			if strings.Contains(strings.ToLower(city), strings.ToLower(search)) { // Check if the city matches
				res = append(res, art)
				break
			}
		}
	}
	res = RemoveDuplicates(res)
	return res, nil
}

// RemoveDuplicates - Remove duplicates from the result list
func RemoveDuplicates(artists []Artist) []Artist {
	seen := make(map[int]bool)
	var uniqueArtists []Artist

	for _, art := range artists {
		if !seen[art.ID] { // Check if the artist's ID is already present, otherwise add it
			seen[art.ID] = true
			uniqueArtists = append(uniqueArtists, art)
		}
	}
	return uniqueArtists
}
