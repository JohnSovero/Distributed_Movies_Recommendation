package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Estructura para representar una película
type Movie struct {
	MovieID  int      `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Genres   []string `json:"genres"`
	IMDBLink string   `json:"imdb_link"`
	TMDBLink string   `json:"tmdb_link"`
}

// Estructura para representar una solicitud de recomendación
type RecommendationRequest struct {
	UserID int `json:"userID"`
	NumRec int `json:"numRec"`
}

// Variables globales para almacenar usuarios y películas
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
	// Endpoint para obtener todas las películas
	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	// Endpoint para obtener todos los usuarios
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	// Endpoint para obtener una película por ID
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	// Endpoint para obtener recomendaciones
	router.HandleFunc("/recommendations/{numRec}/genre/{genre}/users/{id}", getRecommendations).Methods("GET")
	// Endpoint para obtener recomendaciones arriba del promedio usando WebSocket
	// router.HandleFunc("/recommendations/above-average", wsGetAboveAverageRecommendations)

	// Iniciar el servidor en el puerto 9015
	log.Fatal(http.ListenAndServe(":9015", router))
}

// Función principal
func main() {
	loadData()        // Cargar datos iniciales
	defineEndpoints() // Definir los endpoints de la API
}
