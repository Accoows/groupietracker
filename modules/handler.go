package modules

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var API AllData            // Variable globale pour les données de l'API
var tpl *template.Template // Template pour les pages HTML
var err error

type Info struct {
	ArtistID interface{} // Stock les informations de l'artiste concerné
}

func init() {
	c := SafeCounter{values: make(map[string]int)}
	go c.GetAPI() // Récupération des données de l'API toutes les 5 minutes (modifiable si ralentissement)
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
	Relation := Relations{}
	relation := ApiRequest("https://groupietrackers.herokuapp.com/api/relation") // Récupération la relation des artistes
	err = json.Unmarshal(relation, &Relation)                                    // Décode les données de la relation
	if err != nil {
		log.Println(err)
		return
	}
	for i := range API.General.Artists {
		LoadArtistInfos(&API.General.Artists[i], Relation) // Association des relations aux artistes
	}
}

// HomePage - Gestionnaire de la page d'accueil
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandle(http.StatusNotFound, w, "404 not found") // On gère les erreurs de la page
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Status Method Not Allowed")
		return
	}

	API.Incorrect = false

	err = tpl.ExecuteTemplate(w, "homepage.html", API) // Démarrage de la page d'accueil
	if err != nil {
		tpl.ExecuteTemplate(w, "errors.html", http.StatusInternalServerError)
		if err != nil {
			ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
		}
	}
}

// MainPage - Gestionnaire de la page principale
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artists" {
		ErrorHandle(http.StatusNotFound, w, "404 not found")
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Status Method Not Allowed")
		return
	}

	API.Search = API.General
	API.Incorrect = false

	err = tpl.ExecuteTemplate(w, "artistsDisplay.html", API) // Démarrage de la page principale
	if err != nil {
		tpl.ExecuteTemplate(w, "errors.html", http.StatusInternalServerError)
		if err != nil {
			ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
		}
	}
}

// SearchHandler - Gestionnaire de recherche
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		ErrorHandle(http.StatusNotFound, w, "404 not found")
		return
	}
	if r.Method != http.MethodPost {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Method Not Allowed")
		return
	}

	err := r.ParseForm() // On récupère les données du formulaire
	if err != nil {
		log.Println(err)
		return
	}

	search := r.FormValue("search") // On récupère la valeur de la recherche

	if search == "" { // Si la recherche est vide, on affiche tous les artistes
		API.Search = API.General
	} else {
		art, searchErr := API.General.SearchArtist(search) // Sinon, on affiche les artistes correspondants à la recherche
		if searchErr != nil || len(art) == 0 {             // Si la recherche ne correspond à aucun artiste
			API.Incorrect = true   // Dans ce cas, on passe sur l'affichage de la pop up d'erreur en HTML/JS
			API.Search = General{} // On affiche un message d'erreur
		} else {
			API.Incorrect = false
			API.Search = General{Artists: art} // Sinon, on affiche les artistes correspondants
		}
	}

	err = tpl.ExecuteTemplate(w, "artistsDisplay.html", API) // Démarrage de la page de recherche
	if err != nil {
		tpl.ExecuteTemplate(w, "errors.html", http.StatusInternalServerError)
		if err != nil {
			ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
		}
	}
}

// ArtistPage - Gestionnaire de page d'artiste
func ArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/artist/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id, err := strconv.Atoi(r.URL.Path[8:]) // On récupère l'ID de l'artiste et on le convertit en int
	if err != nil {
		ErrorHandle(http.StatusBadRequest, w, err, "400 Bad request")
		return
	}
	if !(id > 0 && id <= len(API.General.Artists)) {
		ErrorHandle(http.StatusNotFound, w, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Method Not Allowed")
		return
	}
	info := &Info{
		ArtistID: API.General.Artists[id-1], // On récupère les informations de l'artiste
	} // On retire 1 pour éviter les erreurs de décalage entre Go et l'ID de l'artiste
	err = tpl.ExecuteTemplate(w, "artistInformations.html", info) // Démarrage de la page de recherche
	if err != nil {
		tpl.ExecuteTemplate(w, "errors.html", http.StatusInternalServerError)
		if err != nil {
			ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
		}
	}
}

// ErrorHandle - Gestion des erreurs du site
func ErrorHandle(ErrorStatus int, w http.ResponseWriter, errC ...interface{}) {
	for _, val := range errC {
		log.Println(val) // ON affiche les erreurs dans la console (à supprimer plus tard pour éviter le spam ou autre)
	}
	w.WriteHeader(ErrorStatus)
	tpl.ExecuteTemplate(w, "errors.html", ErrorStatus)
}
