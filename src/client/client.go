package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Estructura para representar a un usuario
type User struct {
	ID      int
	Ratings map[int]float64
}

// Estructura para enviar similitud al servidor
type SimilarityData struct {
	Similarity float64 `json:"similarity"`
	UserID     string  `json:"userID"`
}

// Estructura para recibir datos del servidor
type ServerData struct {
	MainUserRatings map[int]float64 `json:"user1"`
	OtherUsers      map[int]User    `json:"user2"`
}

// Maneja la conexión del cliente
func handleConnection(conn net.Conn) {
	defer conn.Close()
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo de la conexión:", err)
		return
	}
	input = strings.TrimSpace(input)

	// Deserializar JSON a ServerData
	var serverData ServerData
	json.Unmarshal([]byte(input), &serverData)
	var similarityResults []SimilarityData

	// Calcular similitud coseno entre los usuarios y preparar datos para el servidor
	var wg sync.WaitGroup
	wg.Add(len(serverData.OtherUsers))
	for _, user := range serverData.OtherUsers {
		go func(user User) {
			defer wg.Done()
			similarity := calculateCosineSimilarity(serverData.MainUserRatings, user.Ratings)
			similarityResults = append(similarityResults, SimilarityData{
				Similarity: similarity,
				UserID:     strconv.Itoa(user.ID),
			})
		}(user)
	}
	wg.Wait()

	// Enviar resultados de similitud al servidor
	sendSimilarityResults(similarityResults, conn)
}

// Inicia el servicio de escucha en el puerto especificado
func startListening(port string) {
	address := fmt.Sprintf("localhost:%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error al iniciar el servicio de escucha:", err)
		return
	}
	defer listener.Close()

	for { // Bucle para aceptar conexiones
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// Obtiene entrada del usuario desde el terminal
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Calcula la similitud coseno entre dos conjuntos de valoraciones de usuario
func calculateCosineSimilarity(ratings1 map[int]float64, ratings2 map[int]float64) float64 {
	dotProduct := 0.0
	sumSquares1 := 0.0
	sumSquares2 := 0.0

	for itemID, rating1 := range ratings1 {
		if rating2, exists := ratings2[itemID]; exists {
			dotProduct += rating1 * rating2
			sumSquares1 += rating1 * rating1
			sumSquares2 += rating2 * rating2
		}
	}

	if sumSquares1 == 0 || sumSquares2 == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(sumSquares1) * math.Sqrt(sumSquares2))
}

// Envía los resultados de similitud al servidor
func sendSimilarityResults(similarityData []SimilarityData, conn net.Conn) {
	fmt.Printf("Cantidad de datos enviados: %d \n", len(similarityData))
	jsonData, err := json.Marshal(similarityData)
	if err != nil {
		fmt.Println("Error al serializar datos:", err)
		return
	}
	fmt.Fprintln(conn, string(jsonData))
}

func main() {
	fmt.Print("Ingrese el puerto para iniciar el servicio: ")
	port := getUserInput()
	startListening(port)
}
