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
	UserID int    `json:"userID"`
	NumRec int    `json:"numRec"`
	Genre  string `json:"genre"`
}

var users []int
var movies []Movie

// Configuración del actualizador de WebSocket
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Cambiar esto para mayor seguridad en producción
// 	},
// }

// Función para definir los endpoints de la API
func defineEndpoints() {
	router := mux.NewRouter()

	// http.HandleFunc("/movies/", getAllMovies)
	// http.HandleFunc("/users/", getAllUsers)
	// http.HandleFunc("/movie/", getMovieByID)
	// http.HandleFunc("recommendations/", getRecommendations)

	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	// Endpoint para obtener recomendaciones
	router.HandleFunc("/recommendations/{numRec}/genres/{genre}/users/{id}", getRecommendations).Methods("GET")
	// Endpoint para obtener recomendaciones arriba del promedio usando WebSocket
	// router.HandleFunc("/recommendations/above-average", wsGetAboveAverageRecommendations)

	log.Fatal(http.ListenAndServe(":9015", router))
}

func main() {
	loadData()        // Cargar datos iniciales
	defineEndpoints() // Definir los endpoints de la API
	loadData()
	defineEndpoints()
}
