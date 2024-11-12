package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cambiar esto para mayor seguridad en producción
	},
}

func defineEndpoints() {
	router := mux.NewRouter()

	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	router.HandleFunc("/recommendations/{numRec}/users/{id}", getRecommendationsWS)

	log.Fatal(http.ListenAndServe(":9015", router))
}

func main() {
	loadData()
	defineEndpoints()
}
