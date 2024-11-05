package fc

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Leer los ratings de un archivo CSV
func ReadRatingsFromCSV(filename string) (map[int]User, error) {
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

	userMap := make(map[int]User) // Cambiado a un mapa
	for _, record := range records[1:] { // Saltar el encabezado
		userID, _ := strconv.Atoi(record[0])
		itemID, _ := strconv.Atoi(record[1])
		score, _ := strconv.ParseFloat(record[2], 64)

		if _, exists := userMap[userID]; !exists {
			userMap[userID] = User{ID: userID, Ratings: make(map[int]float64)}
		}
		userMap[userID].Ratings[itemID] = score
	}
	fmt.Println("Users:", len(userMap))
	fmt.Println("Total reviews:", len(records))
	return userMap, nil
}

// Función para dividir los usuarios en 3 grupos
func DivideUsers(users map[int]User, userId int) (map[int]User, map[int]User, map[int]User) {
    group1 := make(map[int]User)
    group2 := make(map[int]User)
    group3 := make(map[int]User)

    groups := []map[int]User{group1, group2, group3}
    currentGroup := 0

    for id, user := range users {
		if user.ID != userId{
			groups[currentGroup][id] = user
        	currentGroup = (currentGroup + 1) % 3
		}
    }
    return group1, group2, group3
}