package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Leer los ratings de un archivo CSV
func ReadRatingsFromCSV(filename string) (map[int]types.User, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	userMap := make(map[int]types.User)         // Cambiado a un mapa
	for _, record := range records[1:] { 		// Saltar el encabezado
		userID, _ := strconv.Atoi(record[0])
		itemID, _ := strconv.Atoi(record[1])
		score, _ := strconv.ParseFloat(record[2], 64)

		if _, exists := userMap[userID]; !exists {
			userMap[userID] = types.User{ID: userID, Ratings: make(map[int]float64)}
		}
		userMap[userID].Ratings[itemID] = score
	}
	fmt.Println("\tUsers:", len(userMap))
	fmt.Println("\tTotal reviews:", len(records))
	return userMap, nil
}

// Función para leer el CSV de películas
func ReadMoviesFromCSV(filename string) (map[int]types.Movie, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	movieMap := make(map[int]types.Movie)
	for _, record := range records[1:] { // Saltar el encabezado
		movieID, _ := strconv.Atoi(record[0])
		title := record[1]
		genres := strings.Split(record[2], "|")

		movieMap[movieID] = types.Movie{
			MovieID: movieID,
			Title:   title,
			Genres:  genres,
		}
	}
	fmt.Println("\tMovies:", len(movieMap))
	return movieMap, nil
}

// Función para dividir los usuarios en 3 grupos
func DivideUsers(users map[int]types.User, userId int, numGroups int) []map[int]types.User {
	groups := make([]map[int]types.User, numGroups)
	for i := range groups {
		groups[i] = make(map[int]types.User)
	}
	currentGroup := 0

	for id, user := range users {
		if user.ID != userId {
			groups[currentGroup][id] = user
			currentGroup = (currentGroup + 1) % numGroups
		}
	}
	return groups
}
