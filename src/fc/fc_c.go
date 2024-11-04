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
    User2 map[int]float64 `json:"user2"`
    ID    string          `json:"id"` // Puede ser un string o int, según lo que necesites
}

var similarities map[int]float64
var mu sync.Mutex
var ch = make(chan int, 3)
var wg = sync.WaitGroup{}

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

func sentToClient(user1 map[int]float64, user2 map[int]float64, idUser string, idClient int) {
    defer wg.Done()
    defer func() { <-ch }()
    for {
		hostClient := clientAddresses[idClient]
        conn, err := net.Dial("tcp", hostClient)
		muLocal := &sync.Mutex{}
        if err == nil {
            defer conn.Close()

            data := ToClientData{
                User1: user1,
                User2: user2,
                ID:    idUser,
            }
            // Serializar la estructura a JSON
            jsonData, err := json.Marshal(data)
            if err != nil {
                fmt.Println("Error al serializar datos:", err)
                return
            }
            // Enviar datos al cliente
            muLocal.Lock()
            _, err = fmt.Fprintln(conn, string(jsonData))
            if err != nil {
                fmt.Println("Error al enviar datos al cliente:", err)
                return
            }
            defer muLocal.Unlock()
            return
        } else{
			muLocal.Lock()
			idClient++
			idClient = idClient % 3
			defer muLocal.Unlock()
		}
        fmt.Printf("Error al conectar al cliente: %v. Reintentando...\n", err)
		fmt.Printf("Intentando con el cliente %d\n", idClient%3)
    }
}

// Función para encontrar los usuarios más similares a un usuario dado
func mostSimilarUsersC(users map[int]User, userID int) []int {
	for id, user := range users {
		if id != userID {
            wg.Add(1)
			ch <- 1 // Limitar a 3 goroutines concurrentes
			go func(user User, id int) {
				sentToClient(users[userID].Ratings, user.Ratings, strconv.Itoa(user.ID), id%3)
			}(user, id)
		}
	}
    wg.Wait()
	// Ordenar los usuarios por similitud
	type kv struct {
		Key   int
		Value float64
	}
	var sortedSimilarities []kv
	mu.Lock()
	for k, v := range similarities {
		sortedSimilarities = append(sortedSimilarities, kv{k, v})
	}
	mu.Unlock()

	sort.Slice(sortedSimilarities, func(i, j int) bool {
		return sortedSimilarities[i].Value > sortedSimilarities[j].Value
	})

	var mostSimilar []int
	for _, kv := range sortedSimilarities {
		mostSimilar = append(mostSimilar, kv.Key)
	}
	return mostSimilar
}

func Handle(con net.Conn) {
    defer func() { 
		//wg.Done()
		//<-ch 
		con.Close()
	}() // Liberar espacio en el canal al finalizar
	msg, err:= bufio.NewReader(con).ReadString('\n')
    msg = strings.TrimSpace(msg)
	if err != nil {
		fmt.Println("Error al leer de la conexión:", err)
		return
	}

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
    // Guardar la similitud en un mapa
	muHandle.Lock()
	fmt.Printf("Recibido: Similitud = %f, ID de Usuario = %d\n", similarity, userIDInt)
	similarities[userIDInt] = similarity
	defer muHandle.Unlock()
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
