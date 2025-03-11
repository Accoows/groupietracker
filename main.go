package main

import (
	"fmt"
	"groupietracker/modules"
	"log"
	"net/http"
)

func main() {
	// Gestion of the static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", modules.HomePage)
	http.HandleFunc("/artists", modules.MainPage)
	http.HandleFunc("/artist/", modules.ArtistPage)
	http.HandleFunc("/search", modules.SearchHandler)

	port := ":8080"
	fmt.Println("Serveur lancé sur http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
