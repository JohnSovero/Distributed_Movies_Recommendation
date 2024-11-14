package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Movie struct {
	MovieID    int      `json:"id"`
	Title      string   `json:"title"`
	Year       int      `json:"year"`
	Genres     []string `json:"genres"`
	IMDBLink   string   `json:"imdb_link"`
	TMDBLink   string   `json:"tmdb_link"`
	Overview   string   `json:"overview"`
	VoteAvg    string   `json:"vote_avg"`
	PosterPath string   `json:"poster"`
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
func defineEndpoints(port string) {
	router := mux.NewRouter()

	// http.HandleFunc("/movies/", getAllMovies)
	// http.HandleFunc("/users/", getAllUsers)
	// http.HandleFunc("/movie/", getMovieByID)
	// http.HandleFunc("recommendations/", getRecommendations)
	port = ":" + port
	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	// Endpoint para obtener recomendaciones
	router.HandleFunc("/recommendations/{numRec}/genres/{genre}/users/{id}", getRecommendations).Methods("GET")
	// Endpoint para obtener recomendaciones arriba del promedio usando WebSocket
	// router.HandleFunc("/recommendations/above-average", wsGetAboveAverageRecommendations)
	router.HandleFunc("/recommendations/above-average", wsGetAboveAverageRecommendations).Methods("GET")

	log.Fatal(http.ListenAndServe(port, router))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9015" // Puerto por defecto si no está configurado
	}
	loadData()
	defineEndpoints(port)
}
