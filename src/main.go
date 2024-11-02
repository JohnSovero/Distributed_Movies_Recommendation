package main

import (
    "log"
	"PC4/fc"
    "PC4/utils"
	"net"

)

func main() {
    // Leer archivo de recomendación de películas
    ratings, err := fc.ReadRatingsFromCSV("dataset/ratings.csv")
    if err != nil {
        log.Fatalf("Error reading ratings from CSV: %v", err)
    }

    // Modo: escucha del servidor
    server := "localhost:9005"
    ln, err := net.Listen("tcp", server)
    if err != nil {
        log.Fatalf("Error al iniciar el servidor: %v", err)
    }
    defer ln.Close()

    user := 1
    // Predecir y manejar conexiones
    utils.PredictFCC(ratings, user, 8, ln)
}
