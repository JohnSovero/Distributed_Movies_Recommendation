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
var clientAddresses = []string{"localhost:8000", "localhost:8001", "localhost:8002"}
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
    for i := 0; i < len(clientAddresses); i++ {
		hostClient := clientAddresses[idClient]
        conn, err := net.Dial("tcp", hostClient)
		if err == nil {
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
            _, err = fmt.Fprintln(conn, string(jsonData))
            if err != nil {
                fmt.Println("Error al enviar datos al cliente:", err)
                return
            }
            // Manejar la conexión del cliente 
			HandleClients(conn)
			return
        } else{
			fmt.Printf("Error al conectar al cliente: %v. Reintentando...\n", err)
			idClient = (idClient + 1 ) % 3
		}
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
	
	// Enviar los datos a los clientes
	wgRecibidos.Add(len(clientAddresses))
	for i, group := range []map[int]User{group1, group2, group3} {
		go func (group map[int]User, i int) {
			sentToClient(users[userID].Ratings, group, i%len(clientAddresses))
		}(group, i)
	}
	wgRecibidos.Wait()
	fmt.Printf("Cantidad de similaridades con usuarios calculadas: %d\n", len(similarities))

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
func HandleClients(con net.Conn) {
	defer wgRecibidos.Done()
	defer con.Close()
	
	reader := bufio.NewReader(con)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer de la conexión:", err)
		return
	}
	msg = strings.TrimSpace(msg)

	// Deserializar JSON a la estructura FromClientData
	var message []FromClientData
	err = json.Unmarshal([]byte(msg), &message)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}
	fmt.Printf("Recibidos %d datos\n", len(message))
	muHandle:= &sync.Mutex{}
	muHandle.Lock()
	for _, data := range message {
		userId, _ := strconv.Atoi(data.UserID)
		similarities[userId] = data.Similarity
	}
	muHandle.Unlock()
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func RecommendItemsC(users map[int]User, userIndex int, numRecommendations int) []int {
	similarities = make(map[int]float64)
	similarUsers := mostSimilarUsersC(users, userIndex)
	recommendations := make(map[int]float64)

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, similarUser := range similarUsers {
		wg.Add(1)
		go func(similarUser int) {
			defer wg.Done()
			for itemID, rating := range users[similarUser].Ratings {
				if _, exists := users[userIndex].Ratings[itemID]; !exists {
					mu.Lock()
					recommendations[itemID] += rating
					mu.Unlock()
				}
			}
		}(similarUser)
	}
	wg.Wait()
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
	//imprimir 10 mejores recomendaciones con similaridad
	//mt.Println("10 mejores recomendaciones con similaridad")
	//or i := 0; i < 10 && i < len(sortedRecommendations); i++ {
	//	fmt.Printf("Item: %d, Similaridad: %f\n", sortedRecommendations[i].Key, sortedRecommendations[i].Value)
	//
	return recommendedItems
}
