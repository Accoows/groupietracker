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
		Relation := Relations{}
		artist := ApiRequest("https://groupietrackers.herokuapp.com/api/artists")
		relation := ApiRequest("https://groupietrackers.herokuapp.com/api/relation")
		err := json.Unmarshal(artist, &API.General.Artists) // Recupération JSON des artistes
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
			LoadArtistInfos(&API.General.Artists[i], Relation) // Décode les données des relations
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

func LoadArtistInfos(art *Artist, relation Relations) {
	for _, val := range relation.Index {
		if art.ID == val.ID {
			newDatesLocations := changeRelationCaracters(val.DatesLocations)
			art.DatesLocations = newDatesLocations // Met à jour les relations si l'ID correspond
		}
	}
}
