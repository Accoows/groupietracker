package main

import (
	"encoding/json"
	"groupietracker/modules" // Import du package models et searchbar
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Variables globales
var artists []modules.Artist
var locations modules.LocationData
var dates modules.DatesData
var relations modules.RelationData

// APIBase récupère les URL des endpoints depuis l'API principale
func APIBase() modules.APIResponse {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		log.Println("Erreur lors de la récupération de l'API :", err)
		return modules.APIResponse{}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse :", err)
		return modules.APIResponse{}
	}

	var api modules.APIResponse
	err = json.Unmarshal(body, &api)
	if err != nil {
		log.Println("Erreur lors de la désérialisation de l'API principale :", err)
		return modules.APIResponse{}
	}
	return api
}

// Charger les données à partir d'un endpoint API
func loadData(url string, target interface{}) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return
	}
}

func displayHomepage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func displayArtists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query") // Requête utilisateur
	var suggestions []modules.Suggestion

	// Si une recherche est effectuée, génère les suggestions
	if query != "" {
		suggestions = modules.Search(query, artists, locations, dates)
	}

	// Structure des données à transmettre au template
	data := struct {
		Query       string
		Suggestions []modules.Suggestion
		Artists     []modules.Artist
	}{
		Query:       query,
		Suggestions: suggestions,
		Artists:     artists,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseFiles("templates/artistsDisplay.html")
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func displayArtistDetails(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	for _, artist := range artists {
		if strconv.Itoa(artist.ID) == id {
			var artistLocations []string
			var artistDates []string
			var artistRelations map[string][]string

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
				Artist    modules.Artist
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
				http.Error(w, "Erreur interne", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
	}
	http.Error(w, "Erreur interne (displayArtistDetails)", http.StatusInternalServerError)
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

	http.HandleFunc("/homepage", displayHomepage)
	http.HandleFunc("/artistsDisplay", displayArtists)

	// Intégration de la recherche via searchbar.go
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		modules.HandleSearch(w, r, artists, locations, dates, relations)
	})

	http.HandleFunc("/artistInformations", displayArtistDetails)
	http.HandleFunc("/", defaultPage)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
