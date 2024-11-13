package model

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
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/types"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/utils"
)

// Variables globales
var clientAddresses = []string{"localhost:8000", "localhost:8001", "localhost:8002"}
var similarityScores map[int]float64
var waitGroupResponses = sync.WaitGroup{}
var mutex = &sync.Mutex{}

const TIMEOUT = 10 * time.Second
const MAX_RETRIES = 2

// 500ms
const RETRY_DELAY = 150 * time.Millisecond

func sentToClient(userRatings map[int]float64, userGroups map[int]types.User, clientID int) {
	var attempts int
	for attempts = 0; attempts < len(clientAddresses); attempts++ {
		clientAddress := clientAddresses[clientID]
		conn, err := net.Dial("tcp", clientAddress)
		if err == nil {
			data := types.ToClientData{
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
		} else {
			fmt.Printf("Error al conectar al cliente: %v. Reintentando...\n", err)
			clientID = (clientID + 1) % len(clientAddresses)
		}
		fmt.Printf("Intentando con el cliente %d\n", clientID%len(clientAddresses))
	}
	if attempts == len(clientAddresses) {
		waitGroupResponses.Done()
		fmt.Println("No hay ningún cliente activo.")
	}
}

// Función para encontrar las similaridades entre un usuario y los demás
func findSimilarUsers(users map[int]types.User, userID int) map[int]float64 {
	mu := &sync.Mutex{}
	mu.Lock()
	groups := utils.DivideUsers(users, userID, len(clientAddresses))
	mu.Unlock()

	// Imprimir la cantidad de usuarios en cada cliente
	fmt.Println("\nDistribución de usuarios por cliente:")
	for i, group := range groups {
		fmt.Printf("\t- Cliente %d: %d\n", i+1, len(group))
	}
	// Inicializar el mapa similarityScores
	similarityScores = make(map[int]float64)
	// Enviar los datos a los clientes
	waitGroupResponses.Add(len(clientAddresses))
	for i, group := range groups {
		go func(group map[int]types.User, i int) {
			sentToClient(users[userID].Ratings, group, i%len(clientAddresses))
		}(group, i)
	}
	waitGroupResponses.Wait()
	fmt.Printf("\nCantidad de similaridades con usuarios calculadas: %d\n", len(similarityScores))

	return similarityScores
}

// Función para manejar las conexiones de los clientes en el servidor
func HandleClients(conn net.Conn) {
	defer waitGroupResponses.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer de la conexión:", err)
		return
	}
	msg = strings.TrimSpace(msg)

	// Deserializar JSON a la estructura FromClientData
	var message []types.FromClientData
	err = json.Unmarshal([]byte(msg), &message)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}
	//fmt.Printf("Recibidos %d datos\n", len(message))

	mutex.Lock()
	defer mutex.Unlock()
	for _, data := range message {
		userId, _ := strconv.Atoi(data.UserID)
		similarityScores[userId] = data.Similarity
	}

}

// Generate recommendations and return the ones with the score above average (i want all, not specifying the number)
func GenerateRecommendationsAboveAverage(users map[int]types.User, userIndex int) []int {
	similarityUsersScores := findSimilarUsers(users, userIndex)
	recommendations := make(map[int]float64)
	averageWeightedRating := 0.0

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for similarUserID, similarity := range similarityUsersScores {
		wg.Add(1)
		go func(similarUserID int, similarity float64) {
			defer wg.Done()

			// Iterar sobre las calificaciones del usuario similar
			for itemID, rating := range users[similarUserID].Ratings {
				if _, exists := users[userIndex].Ratings[itemID]; !exists { // Si el usuario principal no ha calificado el ítem
					mutex.Lock()
					// Ponderamos el rating por la similitud entre el usuario principal y el usuario similar
					weightedRating := rating * similarity
					averageWeightedRating += weightedRating
					recommendations[itemID] += weightedRating
					mutex.Unlock()
				}
			}
		}(similarUserID, similarity)
	}
	wg.Wait()

	// Ordenar las recomendaciones por las calificaciones acumuladas
	var aboveAvgRecs []int
	for k, v := range recommendations {
		if v > averageWeightedRating {
			aboveAvgRecs = append(aboveAvgRecs, k)
		}
	}
	return aboveAvgRecs

	// var sortedRecs []types.Kv
	// for k, v := range recommendations {
	// 	sortedRecs = append(sortedRecs, types.Kv{Key: k, Value: v})
	// }
	// sort.Slice(sortedRecs, func(i, j int) bool {
	// 	return sortedRecs[i].Value > sortedRecs[j].Value
	// })

	// var recommendedItems []int
	// for i := 0; i < len(recommendations); i++ {
	// 	if sortedRecs[i].Value > averageWeightedRating {
	// 		recommendedItems = append(recommendedItems, sortedRecs[i].Key)
	// 	} else {
	// 		break
	// 	}
	// }
	// return recommendedItems
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func GenerateRecommendations(users map[int]types.User, userIndex int, numRecs int) []int {
	similarityUsersScores := findSimilarUsers(users, userIndex)
	recommendations := make(map[int]float64)

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for similarUserID, similarity := range similarityUsersScores {
		wg.Add(1)
		go func(similarUserID int, similarity float64) {
			defer wg.Done()

			// Iterar sobre las calificaciones del usuario similar
			for itemID, rating := range users[similarUserID].Ratings {
				if _, exists := users[userIndex].Ratings[itemID]; !exists { // Si el usuario principal no ha calificado el ítem
					mutex.Lock()
					// Ponderamos el rating por la similitud entre el usuario principal y el usuario similar
					weightedRating := rating * similarity
					recommendations[itemID] += weightedRating
					mutex.Unlock()
				}
			}
		}(similarUserID, similarity)
	}
	wg.Wait()

	// Ordenar las recomendaciones por las calificaciones acumuladas
	var sortedRecs []types.Kv
	for k, v := range recommendations {
		sortedRecs = append(sortedRecs, types.Kv{Key: k, Value: v})
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
func PredictFC(users map[int]types.User, targetUser int, k int, movies map[int]types.Movie) {
	fmt.Printf("\nPredicciones para el usuario %d\n", targetUser)
	start := time.Now()
	recommendationsFCC := GenerateRecommendations(users, targetUser, k)
	elapsed := time.Since(start)

	var movieTitles []string
	for _, movieID := range recommendationsFCC {
		movieTitles = append(movieTitles, movies[movieID].Title)
	}

	fmt.Printf("\nPelículas recomendadas:\n")
	for i, movie := range movieTitles {
		fmt.Printf("\t%d. %s [id: %d]\n", i+1, movie, recommendationsFCC[i])
	}
	fmt.Printf("\nTiempo de ejecución de filtrado colaborativo: %v\n", elapsed)
}
