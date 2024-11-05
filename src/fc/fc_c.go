package fc

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)
var clientAddresses = []string{"localhost:8000", "localhost:8001", "localhost:8002"}
type User struct {
	ID      int
	Ratings map[int]float64
}

type FromClientData struct {
	Similarity float64 `json:"similarity"`
	UserID     string  `json:"userID"`
}
type ToClientData struct {
    User1 map[int]float64 `json:"user1"`
    User2 map[int]User `json:"user2"`
}

type kv struct {
	Key   int
	Value float64
}

var similarities map[int]float64
var wgRecibidos = sync.WaitGroup{}
var userMapQuantity = 0

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
	userMapQuantity = len(userMap)
	return userMap, nil
}

// Función para dividir los usuarios en 3 grupos
func divideUsers(users map[int]User, userId int) (map[int]User, map[int]User, map[int]User) {
    group1 := make(map[int]User)
    group2 := make(map[int]User)
    group3 := make(map[int]User)

    groups := []map[int]User{group1, group2, group3}
    currentGroup := 0

    for id, user := range users {
		if user.ID != userId{
			groups[currentGroup][id] = user
        	currentGroup = (currentGroup + 1) % 3 // Avanza al siguiente grupo en rotación
		}
    }

    return group1, group2, group3
}

func sentToClient(user1 map[int]float64, user2 map[int]User, idClient int) {
    for id := 0; id < len(clientAddresses); id++ {
		hostClient := clientAddresses[idClient]
        conn, err := net.Dial("tcp", hostClient)
		if err == nil {
			defer conn.Close()
            data := ToClientData{
                User1: user1,
                User2: user2,
            }
            // Serializar la estructura a JSON
            jsonData, err := json.Marshal(data)
            if err != nil {
                fmt.Println("Error al serializar datos:", err)
                return
            }
            // Enviar datos al cliente
			mu := &sync.Mutex{}
            mu.Lock()
            _, err = fmt.Fprintln(conn, string(jsonData))
			mu.Unlock()
            if err != nil {
                fmt.Println("Error al enviar datos al cliente:", err)
                return
            }
            return
        } else{
			idClient = (idClient + 1 ) % (len(clientAddresses))
		}
        fmt.Printf("Error al conectar al cliente: %v. Reintentando...\n", err)
		fmt.Printf("Intentando con el cliente %d\n", idClient%len(clientAddresses))
    }
}

// Función para encontrar los usuarios más similares a un usuario dado
func mostSimilarUsersC(users map[int]User, userID int) []int {
	// Esperar que se dividan los usuarios
	mu := &sync.Mutex{}
	mu.Lock()
	group1, group2, group3 := divideUsers(users, userID)
	mu.Unlock()

	wgRecibidos.Add(userMapQuantity-1)
	for i, group := range []map[int]User{group1, group2, group3} {
		sentToClient(users[userID].Ratings, group, i%(len(clientAddresses)))
	}
	wgRecibidos.Wait()

	// Ordenar los usuarios por similitud y devolver los más similares
	var sortedSimilarities []kv
	for k, v := range similarities {
		sortedSimilarities = append(sortedSimilarities, kv{k, v})
	}
	sort.Slice(sortedSimilarities, func(i, j int) bool {
		return sortedSimilarities[i].Value > sortedSimilarities[j].Value
	})
	var mostSimilar []int
	for _, kv := range sortedSimilarities {
		mostSimilar = append(mostSimilar, kv.Key)
	}
	return mostSimilar
}
// Función para manejar las conexiones de los clientes en el servidor
func Handle(con net.Conn) {
    defer func() { 
		wgRecibidos.Done()
		con.Close()
	}()
	reader := bufio.NewReader(con)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer de la conexión:", err)
		return
	}
	msg = strings.TrimSpace(msg)

	// Deserializar JSON a la estructura FromClientData
	var message FromClientData
	err = json.Unmarshal([]byte(msg), &message)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}

	// Acceder a los datos deserializados
	similarity := message.Similarity
	userID := message.UserID

	// Convertir el userID de string a int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Printf("Error al convertir el ID de usuario a entero: %v\n", err)
		return
	}
	muHandle:= &sync.Mutex{}
	muHandle.Lock()
	fmt.Printf("Recibido: Similitud = %f, ID de Usuario = %d\n", similarity, userIDInt)
	similarities[userIDInt] = similarity
	muHandle.Unlock()
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func RecommendItemsC(users map[int]User, userIndex int, numRecommendations int) []int {
	similarities = make(map[int]float64)
	similarUsers := mostSimilarUsersC(users, userIndex)
	recommendations := make(map[int]float64)

	var wg2 sync.WaitGroup
	mu := &sync.Mutex{}

	for _, similarUser := range similarUsers {
		wg2.Add(1)
		go func(similarUser int) {
			defer wg2.Done()
			for itemID, rating := range users[similarUser].Ratings {
				if _, exists := users[userIndex].Ratings[itemID]; !exists {
					mu.Lock()
					recommendations[itemID] += rating
					mu.Unlock()
				}
			}
		}(similarUser)
	}
	wg2.Wait()
    //---------------
	// Ordenar las recomendaciones por las calificaciones acumuladas
	type kv struct {
		Key   int
		Value float64
	}
	var sortedRecommendations []kv
	for k, v := range recommendations {
		sortedRecommendations = append(sortedRecommendations, kv{k, v})
	}
	// Ordenar en orden descendente
	sort.Slice(sortedRecommendations, func(i, j int) bool {
		return sortedRecommendations[i].Value > sortedRecommendations[j].Value
	})
	
	// Devolver los índices de los ítems recomendados
	var recommendedItems []int
	for i := 0; i < numRecommendations && i < len(sortedRecommendations); i++ {
		recommendedItems = append(recommendedItems, sortedRecommendations[i].Key)
	}
	return recommendedItems
}
