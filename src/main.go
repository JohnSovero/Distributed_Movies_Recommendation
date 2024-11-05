package main

import (
    "fmt"
    "log"
	"PC4/fc"
    "PC4/utils"
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
        var user int
        fmt.Print("Ingrese el ID del usuario para obtener recomendaciones (o -1 para salir): ")
        _, err := fmt.Scanf("%d", &user)
        if err != nil || user == -1 {
            fmt.Println("Saliendo...")
            break
        }

        // Predecir y manejar conexiones
        utils.PredictFCC(ratings, user, 8)
    }
}
