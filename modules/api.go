package modules

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// GetAPI - Merge JSON every 5 minutes
func (c *SafeCounter) GetAPI() {
	c.mu.Lock()
	for {
		defer c.mu.Unlock()
		Relation := Relations{}
		artist := ApiRequest("https://groupietrackers.herokuapp.com/api/artists")
		relation := ApiRequest("https://groupietrackers.herokuapp.com/api/relation")
		err := json.Unmarshal(artist, &API.General.Artists) // Retrieve JSON of artists
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(relation, &Relation) // Retrieve JSON of relations
		if err != nil {
			log.Println(err)
			return
		}

		for i := range API.General.Artists {
			LoadArtistInfos(&API.General.Artists[i], Relation) // Decode relation data
		}

		log.Println("API has been updated.")
		time.Sleep(time.Minute * 5)
	}
}

func ApiRequest(url string) []byte { // Retrieve data from the API
	resp, err := http.Get(url) // Retrieve data and store it in a byte array
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func LoadArtistInfos(art *Artist, relation Relations) {
	for _, val := range relation.Index {
		if art.ID == val.ID {
			newDatesLocations := changeRelationCaracters(val.DatesLocations)
			art.DatesLocations = newDatesLocations // Update relations if the ID matches
		}
	}
}

// uniqueCities - Extract all cities where concerts took place
func uniqueCities(artists []Artist) []string {
	citySet := make(map[string]bool)
	for _, artist := range artists {
		for city := range artist.DatesLocations {
			citySet[city] = true
		}
	}

	// Convert the map to a slice
	var cities []string
	for city := range citySet {
		cities = append(cities, city)
	}

	return cities
}
