package modules

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// GetAPI - Merge du JSON toutes les 5 minutes
func (c *SafeCounter) GetAPI() {
	c.mu.Lock()
	for {
		defer c.mu.Unlock()
		Locations := LocationData{}
		Dates := DatesData{}
		Relation := Relations{}
		artist := ApiRequest("https://groupietrackers.herokuapp.com/api/artists")
		location := ApiRequest("https://groupietrackers.herokuapp.com/api/locations")
		date := ApiRequest("https://groupietrackers.herokuapp.com/api/dates")
		relation := ApiRequest("https://groupietrackers.herokuapp.com/api/relation")
		err := json.Unmarshal(artist, &API.General.Artists) // Recupération JSON des artistes
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(location, &Locations) // Recupération JSON des locations
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(date, &Dates) // Recupération JSON des dates
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(relation, &Relation) // Recupération JSON des relations
		if err != nil {
			log.Println(err)
			return
		}

		for i := range API.General.Artists {
			LoadArtistInfos(&API.General.Artists[i], Relation, Locations, Dates) // Décode les données des relations
		}

		log.Println("Api has been updated.")
		time.Sleep(time.Minute * 5)
	}
}

func ApiRequest(url string) []byte { // Récupération des données de l'API
	resp, err := http.Get(url) // On récupère les données et on les stocke dans un tableau de bytes
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

func LoadArtistInfos(art *Artist, relation Relations, locations LocationData, dates DatesData) {
	var artistLocations []string
	var artistDates []string
	for _, loc := range locations.Index {
		if art.ID == loc.ID {
			artistLocations = loc.Locations
			for i := 0; i < len(artistLocations); i++ { // artistLocations is parcoured to separate all the locations
				artistLocations[i] = changeLocationsCaracters(artistLocations, i) //The date in artistLocations is replaced with the location modified
			}
			art.Locations = artistLocations
		}
	}
	for _, date := range dates.Index {
		if art.ID == date.ID {
			artistDates = date.Dates
			for i := 0; i < len(artistDates); i++ { // artistDates is parcoured to separate all the dates
				artistDates[i] = changeDatesCaracters(artistDates, i) //The date in artistDates is replaced with the date modified
			}
			art.Dates = artistDates
		}
	}
	for _, val := range relation.Index {
		if art.ID == val.ID {
			art.DatesLocations = val.DatesLocations // Met à jour les relations si l'ID correspond
		}
	}
}
