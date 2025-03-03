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

// Relations - Structure of concert relations
type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// General structure of data
type General struct {
	Artists []Artist
}

type AllData struct {
	General   General
	Search    General
	Incorrect bool
}

// System to avoid concurrency issues
type SafeCounter struct {
	mu     sync.Mutex
	values map[string]int
}
