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

var artists []Artist

func APIBase() APIResponse {
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

func displayArtistsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du modèle HTML", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artists)
}

func displayArtistDetailsPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	for _, artist := range artists {
		if strconv.Itoa(artist.ID) == id {
			tmpl, err := template.ParseFiles("templates/artist.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement du modèle HTML", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, artist)
			return
		}
	}

	http.Redirect(w, r, "/error", http.StatusFound)
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/artists", http.StatusFound)
}

func main() {
	api := APIBase()
	loadArtists(api.Artists)

	http.HandleFunc("/artists", displayArtistsPage)
	http.HandleFunc("/artist", displayArtistDetailsPage)
	http.HandleFunc("/", defaultPage)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
