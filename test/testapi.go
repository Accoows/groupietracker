package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type APIResponse struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relations string `json:"relation"`
}

type Artist struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Members    []string `json:"members"`
	Creation   int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
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

type RelationData struct {
	Index []struct {
		ID        int                 `json:"id"`
		Relations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var artists []Artist
var locations LocationData
var dates DatesData
var relations RelationData

func main() {
	api := fetchAPIBase()
	loadArtists(api.Artists)
	loadLocations(api.Locations)
	loadDates(api.Dates)
	loadRelations(api.Relations)

	displayAllArtist()
}

func fetchAPIBase() APIResponse {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		log.Fatal("Erreur lors de la récupération de l'API :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture du corps de la réponse :", err)
	}

	var api APIResponse
	err = json.Unmarshal(body, &api)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation de l'API principale :", err)
	}
	return api
}

func loadArtists(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des artistes :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture des artistes :", err)
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation des artistes :", err)
	}
}

func loadLocations(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des localisations :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture des localisations :", err)
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation des localisations :", err)
	}
}

func loadDates(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des dates :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture des dates :", err)
	}

	err = json.Unmarshal(body, &dates)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation des dates :", err)
	}
}

func loadRelations(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur lors de la récupération des relations :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture des relations :", err)
	}

	err = json.Unmarshal(body, &relations)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation des relations :", err)
	}
}

func displayAllArtist() {
	for _, artist := range artists {
		fmt.Printf("-----------------------------------------------\n")
		fmt.Printf("Name : %s\n", artist.Name)
		fmt.Printf("ID : %d\n", artist.ID)
		fmt.Printf("Image : %s\n", artist.Image)

		fmt.Printf("Membres :\n", artist.Members)
		for _, member := range artist.Members {
			fmt.Printf("\t- %s\n", member)
		}

		fmt.Printf("Creation Date : %d\n", artist.Creation)
		fmt.Printf("First Album : %s\n", artist.FirstAlbum)
		// Localisations
		fmt.Printf("Localisations :\n")
		localisationFound := false
		for _, loc := range locations.Index {
			if loc.ID == artist.ID {
				for _, location := range loc.Locations {
					fmt.Printf("\t- %s\n", location)
				}
				localisationFound = true
				break
			}
		}

		if !localisationFound {
			fmt.Println("Localisations : Aucune information disponible")
		}

		// Dates
		fmt.Printf("Concert Dates :\n")
		datesFound := false
		for _, date := range dates.Index {
			if date.ID == artist.ID {
				for _, concertDate := range date.Dates {
					fmt.Printf("\t- %s\n", concertDate)
				}
				break
			}
		}

		if !datesFound {
			fmt.Println("Dates : Aucune information disponible")
		}

		// Relations
		fmt.Printf("Relations :\n")
		relationsFound := false
		for _, relation := range relations.Index {
			if relation.ID == artist.ID {
				for location, relationDates := range relation.Relations {
					fmt.Printf("\t%s:\n", location)
					for _, date := range relationDates {
						fmt.Printf("\t\t- %s\n", date)
					}
				}
				break
			}
		}

		if !relationsFound {
			fmt.Println("Relations : Aucune information disponible")
		}

		fmt.Printf("-----------------------------------------------\n")
	}
}
