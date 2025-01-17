package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Structures pour les données
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

// Variables globales pour stocker les données
var artists []Artist
var locations LocationData
var dates DatesData
var relations RelationData

func main() {
	// Charger toutes les données
	api := fetchAPIBase()
	loadArtists(api.Artists)
	loadLocations(api.Locations)
	loadDates(api.Dates)
	loadRelations(api.Relations)

	// Afficher toutes les informations
	displayAllArtists()
}

// Récupère la base de l'API
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

// Charge les artistes
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

// Charge les localisations
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

// Charge les dates
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

// Charge les relations
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

// Affiche toutes les informations sur tous les artistes
func displayAllArtists() {
	for _, artist := range artists {
		fmt.Printf("\n====================================\n")
		fmt.Printf("Nom : %s\n", artist.Name)
		fmt.Printf("Membres : %v\n", artist.Members)
		fmt.Printf("Création : %d\n", artist.Creation)
		fmt.Printf("Premier album : %s\n", artist.FirstAlbum)

		// Localisations
		localisationFound := false
		for _, loc := range locations.Index {
			if loc.ID == artist.ID {
				fmt.Printf("Localisations : %v\n", loc.Locations)
				localisationFound = true
				break
			}
		}
		if !localisationFound {
			fmt.Println("Localisations : Aucune information disponible")
		}

		// Dates
		datesFound := false
		for _, date := range dates.Index {
			if date.ID == artist.ID {
				fmt.Printf("Dates : %v\n", date.Dates)
				datesFound = true
				break
			}
		}
		if !datesFound {
			fmt.Println("Dates : Aucune information disponible")
		}

		// Relations
		relationsFound := false
		for _, relation := range relations.Index {
			if relation.ID == artist.ID {
				fmt.Printf("Relations (Dates et Localisations) : %v\n", relation.Relations)
				relationsFound = true
				break
			}
		}
		if !relationsFound {
			fmt.Println("Relations : Aucune information disponible")
		}
	}
}
