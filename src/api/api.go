package main

import (
	"log"
	"net/http"
)

type Movie struct {
	MovieID  int      `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Genres   []string `json:"genres"`
	IMDBLink string   `json:"imdb_link"`
	TMDBLink string   `json:"tmdb_link"`
}

var users []int
var movies []Movie

func defineEndpoints() {
	http.HandleFunc("/movies/", getAllMovies)
	http.HandleFunc("/users/", getAllUsers)
	http.HandleFunc("/movie/", getMovieByID)

	log.Fatal(http.ListenAndServe(":9015", nil))
}

func main() {
	loadData()
	defineEndpoints()
}
