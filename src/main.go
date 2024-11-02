package main

import (
    "fmt"
    "log"
	"PC4/fc"
    "PC4/utils"
	"net"
)

func main() {
    // Leer archivo de recomendación de películas
    ratings, err := fc.ReadRatingsFromCSV("dataset/ratings2.csv")
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

    for {
        var user int
        fmt.Print("Ingrese el ID del usuario para obtener recomendaciones (o -1 para salir): ")
        _, err := fmt.Scanf("%d", &user)
        if err != nil || user == -1 {
            fmt.Println("Saliendo...")
            break
        }

        // Predecir y manejar conexiones
        utils.PredictFCC(ratings, user, 8, ln)
    }
}
