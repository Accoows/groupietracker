package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Artist struct { //structure qui permet d'afficher les données contenues dans l'API
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	Image        string          `json:"image"`
	Members      []string        `json:"members"`
	CreationDate int             `json:"creationDate"`
	FirstAlbum   string          `json:"firstAlbum"`
	Relations    json.RawMessage `json:"relations"`
}

type Relation struct { //Structure qui permet d'afficher les lieux et les dates des concerts
	DatesLocations map[string][]string `json:"datesLocations"`
}

var artists []Artist

func main() {
	loadArtists("https://groupietrackers.herokuapp.com/api/artists") //Permet d'utiliser l'API

	http.HandleFunc("/", welcomePage)                //Affiche la première page
	http.HandleFunc("/artists", displayArtists)      //Affiche la page avec tout les groupes
	http.HandleFunc("/artist", displayArtistDetails) //Affiche les pages de descriptions de chaques groupes

	log.Println("Server is running on http://localhost:8080") //Message dans la console
	http.ListenAndServe(":8080", nil)                         //Port de la page
}

func loadArtists(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erreur lors de la récupération de l'API :", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture des données de l'API :", err)
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation des artistes :", err)
	}
}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func displayArtists(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du modèle HTML", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artists)
}

func displayArtistDetails(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	for _, artist := range artists {
		if strconv.Itoa(artist.Id) == id {
			var relations Relation
			if err := json.Unmarshal(artist.Relations, &relations); err != nil {
				log.Printf("Relations pour l'artiste %d non structuré, ignoré : %v", artist.Id, err)
				relations = Relation{DatesLocations: make(map[string][]string)} // Valeur par défaut
			}

			data := struct {
				Artist    Artist
				Relations Relation
			}{
				Artist:    artist,
				Relations: relations,
			}

			tmpl, err := template.ParseFiles("templates/artist.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement du modèle HTML", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
	}
	http.NotFound(w, r)
}
