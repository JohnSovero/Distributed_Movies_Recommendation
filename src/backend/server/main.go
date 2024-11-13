package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/types"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/model"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/utils"
)

type RecommendationRequest struct {
	UserID int `json:"userID"`
	NumRec int `json:"numRec"`
}

func serverStartListening(port string, ratings map[int]types.User) {
	address := "localhost:" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error al iniciar el servicio de escucha:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go serverHandleConnection(conn, ratings)
	}
}

func serverHandleConnection(conn net.Conn, ratings map[int]types.User) {
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

		recommendations := model.GenerateRecommendations(ratings, body.UserID, body.NumRec)

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
	pathRatings := "database/data/ratings25.csv"
	// pathMovies := "dataset/movies25.csv"
	fmt.Println("\nLeyendo archivos de datos...")
	fmt.Println("--------------------------------")
	fmt.Println("Detalle de la información procesada:")
	ratings, err := utils.ReadRatingsFromCSV(pathRatings)
	if err != nil {
		log.Fatalf("Error leyendo los ratings del csv: %v", err)
	}

	fmt.Println("Escuchando")
	serverPort := "9000"
	serverStartListening(serverPort, ratings)
}
