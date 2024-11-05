package main

import (
    "PC4/fc"
	"fmt"
	"log"
)

func main() {
    // Leer archivo de recomendación de películas
    path := "dataset/ratings30.csv"
    ratings, err := fc.ReadRatingsFromCSV(path)
    if err != nil {
        log.Fatalf("Error leyendo los ratings del csv: %v", err)
    }

    // Leer el id usuario al que se le harán recomendaciones
    for {
        var userId int
        fmt.Print("Ingrese el ID del usuario para obtener recomendaciones (o -1 para salir): ")
        _, err := fmt.Scanf("%d", &userId)
        if err != nil || userId == -1 {
            fmt.Println("Saliendo...")
            break
        }
        // Predecir y manejar conexiones
        fc.PredictFC(ratings, userId, 8)
    }
}
