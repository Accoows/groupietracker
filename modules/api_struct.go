package modules

import "sync"

// Artist - Structure of artist data
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

// Relations - Structure of concert relations
type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// Filters - Stores active filters
type Filters struct {
	CD  CreationDate
	FAD FirstAlbumDate
}

// CreationDate - Structure to filter by creation date
type CreationDate struct {
	From string
	To   string
}

// FirstAlbumDate - Structure to filter by first album date (under construction)
type FirstAlbumDate struct {
	From string
	To   string
}

// General structure of data
type General struct {
	Artists []Artist
}

type AllData struct {
	General   General
	Search    General
	Filters   Filters
	Incorrect bool
}

// System to avoid concurrency issues
type SafeCounter struct {
	mu     sync.Mutex
	values map[string]int
}
