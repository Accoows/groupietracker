package modules

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var API AllData            // Global variable for API data
var tpl *template.Template // Template for HTML pages
var err error

type Info struct {
	ArtistID interface{} // Stores information about the concerned artist
}

func init() {
	c := SafeCounter{values: make(map[string]int)}
	go c.GetAPI() // Retrieve API data every 5 minutes (modifiable if slowdown occurs)
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
	Relation := Relations{}
	relation := ApiRequest("https://groupietrackers.herokuapp.com/api/relation") // Retrieve artist relations
	err = json.Unmarshal(relation, &Relation)                                    // Decode relation data
	if err != nil {
		log.Println(err)
		return
	}
	for i := range API.General.Artists {
		LoadArtistInfos(&API.General.Artists[i], Relation) // Associate relations with artists
	}
}

// HomePage - Handler for the home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandle(http.StatusNotFound, w, "404 not found") // Handle page errors
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Status Method Not Allowed")
		return
	}

	API.Incorrect = false

	if err = tpl.ExecuteTemplate(w, "homepage.html", API); err != nil {
		ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
	}
}

// MainPage - Handler for the main page
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

	if err := tpl.ExecuteTemplate(w, "artistsDisplay.html", API); err != nil {
		ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
	}
}

// SearchHandler - Handler for search
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		ErrorHandle(http.StatusNotFound, w, "404 not found")
		return
	}
	if r.Method != http.MethodPost {
		ErrorHandle(http.StatusMethodNotAllowed, w, "405 Method Not Allowed")
		return
	}

	err := r.ParseForm() // Retrieve form data
	if err != nil {
		log.Println(err)
		return
	}

	search := r.FormValue("search") // Retrieve search value

	if search == "" { // If search is empty, display all artists
		API.Search = API.General
	} else {
		art, searchErr := API.General.SearchArtist(search) // Otherwise, display artists matching the search
		if searchErr != nil || len(art) == 0 {             // If search does not match any artist
			API.Incorrect = true   // In this case, display the error pop-up in HTML/JS
			API.Search = General{} // Display an error message
		} else {
			API.Incorrect = false
			API.Search = General{Artists: art} // Otherwise, display matching artists
		}
	}

	if err = tpl.ExecuteTemplate(w, "artistsDisplay.html", API); err != nil {
		ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
	}
}

// ArtistPage - Handler for artist page
func ArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/artist/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id, err := strconv.Atoi(r.URL.Path[8:]) // Retrieve artist ID and convert to int
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
		ArtistID: API.General.Artists[id-1], // Retrieve artist information
	} // Subtract 1 to avoid offset errors between Go and artist ID

	if err = tpl.ExecuteTemplate(w, "artistInformations.html", info); err != nil {
		ErrorHandle(http.StatusInternalServerError, w, err, "500 Internal Server Error")
	}
}

// ErrorHandle - Handle site errors
func ErrorHandle(ErrorStatus int, w http.ResponseWriter, errC ...interface{}) {
	for _, val := range errC {
		log.Println(val) // Display errors in the console (to be removed later to avoid spam or other issues)
	}
	w.WriteHeader(ErrorStatus)
	tpl.ExecuteTemplate(w, "errors.html", ErrorStatus)
}
