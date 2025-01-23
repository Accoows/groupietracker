package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
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

func APIBase() APIResponse {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		log.Println("Erreur lors de la récupération de l'API :", err)
		return APIResponse{}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse :", err)
		return APIResponse{}
	}
	var api APIResponse
	err = json.Unmarshal(body, &api)
	if err != nil {
		log.Println("Erreur lors de la désérialisation de l'API principale :", err)
		return APIResponse{}
	}
	return api
}

func loadData(url string, target interface{}) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Erreur lors de la récupération des données depuis %s : %v", url, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Erreur lors de la lecture des données depuis %s : %v", url, err)
		return
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		log.Printf("Erreur lors de la désérialisation des données depuis %s : %v", url, err)
		return
	}
}

func displayWelcome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Println("Erreur lors du chargement du modèle HTML de la page d'accueil")
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func displayArtists(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/artistsDisplay.html")
	if err != nil {
		log.Println("Erreur lors du chargement du modèle HTML")
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artists)
}

func displayArtistDetails(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	for _, artist := range artists {
		if strconv.Itoa(artist.ID) == id {
			var artistLocations []string
			var artistDates []string
			var artistRelations map[string][]string

			// Extraire les données de localisations, dates et relations
			for _, loc := range locations.Index {
				if loc.ID == artist.ID {
					artistLocations = loc.Locations
					break
				}
			}
			for _, date := range dates.Index {
				if date.ID == artist.ID {
					artistDates = date.Dates
					break
				}
			}
			for _, rel := range relations.Index {
				if rel.ID == artist.ID {
					artistRelations = rel.Relations
					break
				}
			}

			data := struct {
				Artist    Artist
				Locations []string
				Dates     []string
				Relations map[string][]string
			}{
				Artist:    artist,
				Locations: artistLocations,
				Dates:     artistDates,
				Relations: artistRelations,
			}

			tmpl, err := template.ParseFiles("templates/artistInformations.html")
			if err != nil {
				log.Println("Erreur lors du chargement du modèle HTML")
				http.Error(w, "Erreur interne", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
	}
	http.Redirect(w, r, "/error", http.StatusFound)
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/homepage", http.StatusFound)
}

func main() {
	api := APIBase()
	loadData(api.Artists, &artists)
	loadData(api.Locations, &locations)
	loadData(api.Dates, &dates)
	loadData(api.Relations, &relations)

	fs := http.FileServer(http.Dir("styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	http.HandleFunc("/homepage", displayWelcome)
	http.HandleFunc("/artistsDisplay", displayArtists)
	http.HandleFunc("/artistInformations", displayArtistDetails)
	http.HandleFunc("/", defaultPage)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
