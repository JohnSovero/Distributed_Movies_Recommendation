package fc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Variables globales
var clientAddresses = []string{"localhost:8000", "localhost:8001", "localhost:8002"}
var similarityScores map[int]float64
var waitGroupResponses = sync.WaitGroup{}

func sentToClient(userRatings map[int]float64, userGroups map[int]User, clientID int) {
    for attempt  := 0; attempt  < len(clientAddresses); attempt ++ {
		clientAddress := clientAddresses[clientID]
        conn, err := net.Dial("tcp", clientAddress)
		if err == nil {
            data := ToClientData{
                User1: userRatings,
                User2: userGroups,
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
			clientID = (clientID + 1) % len(clientAddresses)
		}
		fmt.Printf("Intentando con el cliente %d\n", clientID%len(clientAddresses))
    }
}

// Función para encontrar los usuarios más similares a un usuario dado
func findMostSimilarUsers(users map[int]User, userID int) []int {
	mu := &sync.Mutex{}
	mu.Lock()
	group1, group2, group3 := DivideUsers(users, userID)
	mu.Unlock()
	
	// Enviar los datos a los clientes
	waitGroupResponses.Add(len(clientAddresses))
	for i, group := range []map[int]User{group1, group2, group3} {
		go func (group map[int]User, i int) {
			sentToClient(users[userID].Ratings, group, i%len(clientAddresses))
		}(group, i)
	}
	waitGroupResponses.Wait()
	fmt.Printf("Cantidad de similaridades con usuarios calculadas: %d\n", len(similarityScores))

	// Ordenar los usuarios por similitud y devolver los más similares
	var sortedSimilarities []kv
	for k, v := range similarityScores {
		sortedSimilarities = append(sortedSimilarities, kv{k, v})
	}
	sort.Slice(sortedSimilarities, func(i, j int) bool {
		return sortedSimilarities[i].Value > sortedSimilarities[j].Value
	})
	var mostSimilar []int
	for _, pair := range sortedSimilarities {
		mostSimilar = append(mostSimilar, pair.Key)
	}
	return mostSimilar
}
// Función para manejar las conexiones de los clientes en el servidor
func HandleClients(conn net.Conn) {
	defer waitGroupResponses .Done()
	defer conn.Close()
	
	reader := bufio.NewReader(conn)
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

	var mutex = &sync.Mutex{}
	mutex.Lock()
	for _, data := range message {
		userId, _ := strconv.Atoi(data.UserID)
		similarityScores[userId] = data.Similarity
	}
	mutex.Unlock()
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func generateRecommendations(users map[int]User, userIndex int, numRecs int) []int {
	similarityScores = make(map[int]float64)
	similarUsers := findMostSimilarUsers(users, userIndex)
	recommendations := make(map[int]float64)

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for _, similarUser := range similarUsers {
		wg.Add(1)
		go func(similarUser int) {
			defer wg.Done()
			for itemID, rating := range users[similarUser].Ratings {
				if _, exists := users[userIndex].Ratings[itemID]; !exists {
					mutex.Lock()
					recommendations[itemID] += rating
					mutex.Unlock()
				}
			}
		}(similarUser)
	}
	wg.Wait()

	// Ordenar las recomendaciones por las calificaciones acumuladas
	var sortedRecs []kv
	for k, v := range recommendations {
		sortedRecs = append(sortedRecs, kv{k, v})
	}
	sort.Slice(sortedRecs, func(i, j int) bool {
		return sortedRecs[i].Value > sortedRecs[j].Value
	})
	
	var recommendedItems []int
	for i := 0; i < numRecs && i < len(sortedRecs); i++ {
		recommendedItems = append(recommendedItems, sortedRecs[i].Key)
	}
	return recommendedItems
}

// Recomienda películas a un usuario objetivo utilizando filtrado colaborativo e indica el tiempo de ejecución
func PredictFC(users map[int]User, targetUser int, k int) {
	start := time.Now()
	fmt.Printf("Predicciones para el usuario %d\n", targetUser)
	recommendationsFCC := generateRecommendations(users, targetUser, k)
	fmt.Printf("Recomendaciones de filtrado colaborativo distribuido: %v\n", recommendationsFCC)
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución de filtrado colaborativo: %v\n", elapsed)
}