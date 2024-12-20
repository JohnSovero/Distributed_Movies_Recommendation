package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/types"
	"github.com/JohnSovero/Distributed_Movies_Recommendation/src/backend/utils"
)

// Variables globales
var clientAddresses = []string{
	os.Getenv("NODO1"),
	os.Getenv("NODO2"),
	os.Getenv("NODO3"),
	os.Getenv("NODO4"),
	os.Getenv("NODO5"),
	os.Getenv("NODO6"),
	os.Getenv("NODO7"),
	os.Getenv("NODO8"),
	os.Getenv("NODO9"),
}

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
		conn, err := net.DialTimeout("tcp", clientAddress, TIMEOUT)
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
			fmt.Printf("Error al conectar al cliente: %v. Reintentando... con el client %d\n", err, (clientID+1)%len(clientAddresses))
			clientID = (clientID + 1) % len(clientAddresses)
		}
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
func GenerateRecommendationsAboveAverage(users map[int]types.User, userIndex int, movies map[int]types.Movie, numRecs int) []types.Movie {
	similarityUsersScores := findSimilarUsers(users, userIndex)
	recommendations := make(map[int]float64)
	sumWeightedRating := 0.0

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
					sumWeightedRating += weightedRating
					recommendations[itemID] += weightedRating
					mutex.Unlock()
				}
			}
		}(similarUserID, similarity)
	}
	wg.Wait()

	// compute average rating
	averageRating := sumWeightedRating / float64(len(recommendations))

	var wg2 sync.WaitGroup
	var mu sync.Mutex
	var recommendedMovies []types.Movie

	// get the movies with the weighted rating above average
	for movieID, rating := range recommendations {
		wg2.Add(1)
		go func(movieID int, rating float64) {
			defer wg2.Done()
			if rating > averageRating {
				mu.Lock() // Lock before modifying shared resource
				recommendedMovies = append(recommendedMovies, movies[movieID])
				mu.Unlock() // Unlock after modifying shared resource
			}
		}(movieID, rating)
	}
	wg2.Wait() // Wait for all goroutines to finish

	return recommendedMovies[:numRecs]
}

// Función para recomendar ítems a un usuario basado en usuarios similares
func GenerateRecommendationsByGenre(users map[int]types.User, userIndex int, numRecs int, movies map[int]types.Movie, genre string) []types.Movie {
	similarityUsersScores := findSimilarUsers(users, userIndex)
	recommendations := make(map[int]float64)
	sumWeightedRating := 0.0

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
					sumWeightedRating += weightedRating
					recommendations[itemID] += weightedRating
					mutex.Unlock()
				}
			}
		}(similarUserID, similarity)
	}
	wg.Wait()

	// compute average rating
	averageRating := sumWeightedRating / float64(len(recommendations))

	var wg2 sync.WaitGroup
	var mu sync.Mutex
	var moviesGenre []types.Movie

	// get the movies with the genre and the weighted rating above average
	for movieID, rating := range recommendations {
		wg2.Add(1)
		go func(movieID int, rating float64) {
			defer wg2.Done()
			for _, movie := range movies {
				if movie.MovieID == movieID {
					// Check if the genre exists in the movie's genres list
					for _, g := range movie.Genres {
						if g == genre && rating > averageRating {
							mu.Lock() // Lock before modifying shared resource
							moviesGenre = append(moviesGenre, movie)
							mu.Unlock() // Unlock after modifying shared resource
							break       // No need to check other genres if we found the match
						}
					}
				}
			}
		}(movieID, rating)
	}
	wg2.Wait() // Wait for all goroutines to finish

	return moviesGenre[:numRecs]
}

func GenerateTopRecommendations(users map[int]types.User, userIndex int, movies map[int]types.Movie, numRecs int) []types.Movie {
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

	// Convert the recommendations map to a slice of pairs
	type recommendationPair struct {
		MovieID int
		Rating  float64
	}
	var recommendationPairs []recommendationPair
	for movieID, rating := range recommendations {
		recommendationPairs = append(recommendationPairs, recommendationPair{MovieID: movieID, Rating: rating})
	}

	// Sort the recommendation pairs by rating in descending order
	sort.Slice(recommendationPairs, func(i, j int) bool {
		return recommendationPairs[i].Rating > recommendationPairs[j].Rating
	})

	// Get the top numRecs recommendations
	var recommendedMovies []types.Movie
	for i := 0; i < numRecs && i < len(recommendationPairs); i++ {
		recommendedMovies = append(recommendedMovies, movies[recommendationPairs[i].MovieID])
	}

	return recommendedMovies
}
