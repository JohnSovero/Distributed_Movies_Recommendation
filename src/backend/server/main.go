package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/server/model"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/types"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/utils"
)

type RecommendationRequest struct {
	UserID int    `json:"userID"`
	NumRec int    `json:"numRec"`
	Genre  string `json:"genre"`
}

func serverStartListening(port string, ratings map[int]types.User, name string, movies map[int]types.Movie) {
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	log.Println("Server listening on", address)
	if err != nil {
		log.Println("Error al iniciar el servicio de escucha:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server %s escuchando en el puerto %s\n", name, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go serverHandleConnection(conn, ratings, movies)
	}
}

func serverHandleConnection(conn net.Conn, ratings map[int]types.User, movies map[int]types.Movie) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Conexión cerrada por el cliente.")
				return
			}
			log.Println("Error al leer mensaje:", err)
			return
		}

		var body RecommendationRequest
		json.Unmarshal([]byte(message), &body)

		fmt.Println("Mensaje recibido:", body)
		fmt.Println("Genre received", body.Genre)

		var recommendations []types.Movie

		if body.Genre == "All" {
			recommendations = model.GenerateRecommendationsAboveAverage(ratings, body.UserID, movies, body.NumRec)
		} else {
			recommendations = model.GenerateRecommendationsByGenre(ratings, body.UserID, body.NumRec, movies, body.Genre)
		}

		recommendationsJSON, err := json.Marshal(recommendations)
		if err != nil {
			log.Println("Error al serializar recomendaciones:", err)
			return
		}
		fmt.Fprintln(conn, string(recommendationsJSON))
	}
}

func main() {
	// Leer archivo de recomendación de películas
	pathRatings := "ratings_complete.csv"
	pathMovies := "movies_complete.csv"
	fmt.Println("\nLeyendo archivos de datos...")
	fmt.Println("--------------------------------")
	fmt.Println("Detalle de la información procesada:")
	ratings, err := utils.ReadRatingsFromCSV(pathRatings)
	if err != nil {
		log.Fatalf("Error leyendo los ratings del csv: %v", err)
	}
	movies, err := utils.ReadMoviesFromCSV(pathMovies)
	if err != nil {
		log.Fatalf("Error leyendo las películas del csv: %v", err)
	}

	fmt.Println("Escuchando")
	serverPort := os.Getenv("PORT")

	name := os.Getenv("NODE_NAME")
	if serverPort == "" {
		log.Fatal("El puerto no está configurado en la variable de entorno PORT")
	}
	serverStartListening(serverPort, ratings, name, movies)
}
