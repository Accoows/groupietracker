package modules

import (
	"strconv"
	"strings"
)

// SearchArtist - Recherche un artiste en fonction d'un mot-clé
func (artistsData *General) SearchArtist(search string) ([]Artist, error) {
	var res []Artist

	// Recherche par date de création
	a, atoiErr := strconv.Atoi(search)
	if atoiErr == nil { // Recherche réussie
		for _, art := range artistsData.Artists {
			if art.CreationDate == a { // Parcours les artistes et compare les dates de création
				res = append(res, art)
			}
		}
	}

	// Recherche par nom de l’artiste
	for _, art := range artistsData.Artists {
		if strings.Contains(strings.ToLower(art.Name), strings.ToLower(search)) {
			res = append(res, art)
			continue
		}

		// Recherche par premier album
		if strings.Contains(strings.ToLower(art.FirstAlbum), strings.ToLower(search)) {
			res = append(res, art)
			continue
		}

		// Recherche par membres
		for _, member := range art.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(search)) {
				res = append(res, art)
				break
			}
		}

		// Recherche par dates de concerts dans Relations (DatesLocations)
		for _, dates := range art.DatesLocations {
			for _, date := range dates {
				if strings.Contains(date, search) { // Vérifie si la date saisie correspond
					res = append(res, art)
					break
				}
			}
		}
	}

	// Supprimer les doublons
	res = RemoveDuplicates(res)
	return res, nil
}

// RemoveDuplicates - Supprime les doublons dans la liste des résultats
func RemoveDuplicates(artists []Artist) []Artist {
	seen := make(map[int]bool)
	var uniqueArtists []Artist

	for _, art := range artists {
		if !seen[art.ID] { // Vérifie si l'ID de l'artiste est déjà présent sinon il le rajoute
			seen[art.ID] = true
			uniqueArtists = append(uniqueArtists, art)
		}
	}
	return uniqueArtists
}
