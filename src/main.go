package main

import (
    "fmt"
    "log"
	"PC4/fc"
    "PC4/utils"
)


func main() {
    // Leer archivo de recomendación de películas
    ratings, err := fc.ReadRatingsFromCSV("dataset/ml-32m/ratings30.csv")
    if err != nil {
        log.Fatalf("Error reading ratings from CSV: %v", err)
    }

    //Rol servidor, modo escucha
    go utils.ServicioEscuchar()

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
