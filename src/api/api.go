package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	MovieID  int      `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Genres   []string `json:"genres"`
	IMDBLink string   `json:"imdb_link"`
	TMDBLink string   `json:"tmdb_link"`
}

type RecommendationRequest struct {
	UserID int `json:"userID"`
	NumRec int `json:"numRec"`
}

var users []int
var movies []Movie

func defineEndpoints() {
	router := mux.NewRouter()

	// http.HandleFunc("/movies/", getAllMovies)
	// http.HandleFunc("/users/", getAllUsers)
	// http.HandleFunc("/movie/", getMovieByID)
	// http.HandleFunc("recommendations/", getRecommendations)

	router.HandleFunc("/movies/", getAllMovies).Methods("GET")
	router.HandleFunc("/users/", getAllUsers).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	router.HandleFunc("/recommendations/{numRec}/users/{id}", getRecommendations).Methods("GET")

	log.Fatal(http.ListenAndServe(":9015", router))
}

func main() {
	loadData()
	defineEndpoints()
}
