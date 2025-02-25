package modules

import (
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Suggestion structure les suggestions de recherche
type Suggestion struct {
	Name string
	Type string
}

// HandleSearch gère les requêtes de recherche
func HandleSearch(w http.ResponseWriter, r *http.Request, artists []Artist, locations LocationData, dates DatesData, relations RelationData) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("")))
	if query == "" {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	suggestions := Search(query, artists, locations, dates)

	tmpl, err := template.ParseFiles("templates/artistsDisplay.html")
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		Query       string
		Suggestions []Suggestion
	}{
		Query:       query,
		Suggestions: suggestions,
	})
}

// Recherche principale
func Search(query string, artists []Artist, locations LocationData, dates DatesData) []Suggestion {
	var results []Suggestion

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, Suggestion{Name: artist.Name, Type: "Artiste/Band"})
		}
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				results = append(results, Suggestion{Name: member, Type: "Membre de " + artist.Name})
			}
		}
	}

	// Recherche dans les localisations
	for _, loc := range locations.Index {
		for _, location := range loc.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				results = append(results, Suggestion{Name: location, Type: "location"})
			}
		}
	}

	// Recherche dans les dates de concerts
	for _, dateIndex := range dates.Index {
		for _, date := range dateIndex.Dates {
			if strings.Contains(strings.ToLower(date), query) {
				results = append(results, Suggestion{Name: date, Type: "concert date"})
			}
		}
	}

	// Recherche dans les dates de création et premiers albums
	for _, artist := range artists {
		if strings.Contains(strconv.Itoa(artist.Creation), query) {
			results = append(results, Suggestion{Name: strconv.Itoa(artist.Creation), Type: "creation date of " + artist.Name})
		}
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			results = append(results, Suggestion{Name: artist.FirstAlbum, Type: "first album of " + artist.Name})
		}
	}

	// Trier les résultats par ordre alphabétique
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}
