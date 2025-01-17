package main

type Artist struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Members      string `json:"members"`
	CreationDate int    `json:"creationdate"`
	FisrtAlbum   string `json:"firstalbum"`
}

type Locations struct {
	ID        int    `json:"id"`
	Locations string `json:"locations"`
}

type dates struct {
	ID    int    `json:"id"`
	Dates string `json:"dates"`
}

type relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"dateslocations"`
}

func main() {

}
