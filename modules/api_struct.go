package modules

import "sync"

// Artist - Structure des données des artistes
type Artist struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Image          string   `json:"image"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	DatesLocations map[string][]string
}

type LocationData struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type DatesData struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

// Relations - Structure des relations de concerts
type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// Structure général des données
type General struct {
	Artists   []Artist
	Dates     []DatesData
	Locations []LocationData
}

type AllData struct {
	General   General
	Search    General
	Incorrect bool
}

// Système pour évincer les problèmes de concurrence
type SafeCounter struct {
	mu     sync.Mutex
	values map[string]int
}
