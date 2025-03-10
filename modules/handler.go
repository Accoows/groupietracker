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

	// Add the list of unique cities to API
	API.Filters.City = uniqueCities(API.General.Artists)

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

	// Retrieve form values
	search := r.FormValue("search")             // Search by keyword
	fromCreation := r.FormValue("fromCreation") // Filter: creation date (start)
	toCreation := r.FormValue("toCreation")     // Filter: creation date (end)
	fromFAD := r.FormValue("fromFAD")           // Filter: first album date (start)
	toFAD := r.FormValue("toFAD")               // Filter: first album date (end)
	fromNBOM := r.FormValue("fromNBOM")         // Filter: minimum number of members
	toNBOM := r.FormValue("toNBOM")             // Filter: maximum number of members
	selectedCities := r.Form["cities"]          // Retrieve selected cities

	// Convert dates to numbers
	fromCDYear, errFromCD := strconv.Atoi(fromCreation)
	toCDYear, errToCD := strconv.Atoi(toCreation)
	fromFADYear, errFromFAD := strconv.Atoi(fromFAD)
	toFADYear, errToFAD := strconv.Atoi(toFAD)
	fromNBOMVal, errFromNBOM := strconv.Atoi(fromNBOM)
	toNBOMVal, errToNBOM := strconv.Atoi(toNBOM)

	if errFromCD != nil {
		fromCDYear = 1958
	}
	if errToCD != nil {
		toCDYear = 2024
	}
	if errFromFAD != nil {
		fromFADYear = 1958
	}
	if errToFAD != nil {
		toFADYear = 2024
	}
	// Default values if empty
	if errFromNBOM != nil {
		fromNBOMVal = 1 // Minimum 1 member
	}
	if errToNBOM != nil {
		toNBOMVal = 8 // High value to include everyone
	}

	log.Println("[CD] Filtering - Years:", fromCDYear, toCDYear)
	log.Println("[FAD] Filtering - Years:", fromFADYear, toFADYear)
	log.Println("[NBOM] Filtering - Number:", fromNBOM, toNBOMVal)

	// Initialize filtered artists
	var filteredArtists []Artist

	// Case where no filter is activated: keep all artists
	filteredArtists = API.General.Artists

	// Apply creation date filter if a value is provided
	if fromCDYear > 1958 || toCDYear > 2024 {
		var tempArtists []Artist
		for _, artist := range filteredArtists {
			if artist.CreationDate >= fromCDYear && artist.CreationDate <= toCDYear {
				tempArtists = append(tempArtists, artist)
			}
		}
		filteredArtists = tempArtists // Update the list with the applied filter
	}

	// Apply first album date filter if a value is provided
	if fromFADYear > 1958 || toFADYear > 2024 {
		var tempArtists []Artist
		for _, artist := range filteredArtists {
			albumYear, err := strconv.Atoi(artist.FirstAlbum[len(artist.FirstAlbum)-4:]) // Extract the year
			if err == nil && albumYear >= fromFADYear && albumYear <= toFADYear {
				tempArtists = append(tempArtists, artist)
			}
		}
		filteredArtists = tempArtists // Update the list with the applied filter
	}

	// Apply the number of members filter if a value is provided
	if fromNBOMVal > 0 || toNBOMVal > 8 {
		var tempArtists []Artist
		for _, artist := range filteredArtists {
			numMembers := len(artist.Members)
			if numMembers >= fromNBOMVal && numMembers <= toNBOMVal {
				tempArtists = append(tempArtists, artist)
			}
		}
		filteredArtists = tempArtists // Update the filtered list
	}

	// Filter by concert city
	if len(selectedCities) > 0 {
		var tempArtists []Artist
		for _, artist := range filteredArtists {
			for _, city := range selectedCities {
				if _, exists := artist.DatesLocations[city]; exists {
					tempArtists = append(tempArtists, artist)
					break
				}
			}
		}
		filteredArtists = tempArtists
	}

	log.Println("Artists after filtering (Creation + First Album):", len(filteredArtists))

	// Apply keyword search if a term is entered
	if search != "" {
		general := General{Artists: filteredArtists}
		filteredArtists, err = general.SearchArtist(search)

		// If no artist matches
		if err != nil || len(filteredArtists) == 0 {
			API.Incorrect = true   // Activate error pop-up
			API.Search = General{} // Clear results
		} else {
			API.Incorrect = false
			API.Search = General{Artists: filteredArtists} // Update filtered results
		}
	} else {
		API.Search = General{Artists: filteredArtists} // If no search, keep filtered results or display all artists
	}

	// Gestion des villes pour garder la liste complète
	allCities := uniqueCities(API.General.Artists) // Récupérer toutes les villes
	selectedCityMap := make(map[string]bool)
	for _, city := range selectedCities {
		selectedCityMap[city] = true
	}

	// Organiser les villes : d'abord les cochées, puis le reste
	var updatedCities []string
	for _, city := range allCities {
		if selectedCityMap[city] {
			updatedCities = append(updatedCities, city) // Mettre en haut celles sélectionnées
		}
	}
	for _, city := range allCities {
		if !selectedCityMap[city] {
			updatedCities = append(updatedCities, city)
		}
	}

	API.Filters.City = updatedCities // Mettre à jour l'affichage des villes

	// Save filters for display in the form
	API.Filters.CD.From = fromCreation
	API.Filters.CD.To = toCreation
	API.Filters.FAD.From = fromFAD
	API.Filters.FAD.To = toFAD
	API.Filters.NBOM.From = fromNBOM
	API.Filters.NBOM.To = toNBOM

	// Handle error display
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
