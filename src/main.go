package main

import (
    "PC4/fc"
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func main() {
    // Leer archivo de recomendación de películas
    pathRatings := "dataset/ratings25.csv"
    pathMovies := "dataset/movies25.csv"
    ratings, err := fc.ReadRatingsFromCSV(pathRatings)
    if err != nil {
        log.Fatalf("Error leyendo los ratings del csv: %v", err)
    }
    movies, err := fc.ReadMoviesFromCSV(pathMovies)
    if err != nil {
        log.Fatalf("Error leyendo las películas del csv: %v", err)
    }

    var userId int
    var numRecommendations int = 5 // Valor por defecto
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("\nMenú de opciones:")
        fmt.Println("1. Ingresar ID del usuario")
        fmt.Println("2. Predecir recomendaciones")
        fmt.Println("3. Indicar cuántas películas recomendar")
        fmt.Println("4. Salir")
        fmt.Print("Seleccione una opción: ")

        optionStr, _ := reader.ReadString('\n')
        optionStr = strings.TrimSpace(optionStr)
        option, err := strconv.Atoi(optionStr)
        if err != nil {
            fmt.Println("Opción inválida. Intente de nuevo.")
            continue
        }

        switch option {
        case 1:
            fmt.Print("Ingrese el ID del usuario: ")
            userIdStr, _ := reader.ReadString('\n')
            userIdStr = strings.TrimSpace(userIdStr)
            userId, err = strconv.Atoi(userIdStr)
            if err != nil {
                fmt.Println("ID de usuario inválido. Intente de nuevo.")
                userId = -1
            }
        case 2:
            if userId == -1 {
                fmt.Println("Primero debe ingresar un ID de usuario válido.")
            } else {
                fc.PredictFC(ratings, userId, numRecommendations, movies)
            }
        case 3:
            fmt.Print("Ingrese el número de películas a recomendar: ")
            numRecsStr, _ := reader.ReadString('\n')
            numRecsStr = strings.TrimSpace(numRecsStr)
            numRecommendations, err = strconv.Atoi(numRecsStr)
            if err != nil || numRecommendations <= 0 {
                fmt.Println("Número de recomendaciones inválido. Usando valor por defecto (5).")
                numRecommendations = 5
            }
        case 4:
            fmt.Println("Saliendo...")
            return
        default:
            fmt.Println("Opción inválida. Intente de nuevo.")
        }
    }
}