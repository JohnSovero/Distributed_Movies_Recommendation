package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("Error leyendo de la conexión:", err)
		return
	}
	input = strings.TrimSpace(input)

	// Deserializar JSON a ServerData
	var serverData ServerData
	err = json.Unmarshal([]byte(input), &serverData)
	if err != nil {
		log.Println("Error deserializando datos:", err)
		return
	}
	var similarityResults []SimilarityData

	// Calcular similitud coseno entre los usuarios y preparar datos para el servidor
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(serverData.OtherUsers))
	for _, user := range serverData.OtherUsers {
		go func(user User) {
			defer wg.Done()
			// fmt.Printf("Calculando similitud para usuario %d\n", user.ID)
			similarity := calculateCosineSimilarity(serverData.MainUserRatings, user.Ratings)
			mu.Lock()
			similarityResults = append(similarityResults, SimilarityData{
				Similarity: similarity,
				UserID:     strconv.Itoa(user.ID),
			})
			mu.Unlock()
		}(user)
	}
	wg.Wait()
	// Enviar resultados de similitud al servidor
	sendSimilarityResults(similarityResults, conn)
}

// Inicia el servicio de escucha en el puerto especificado
func startListening(port string, name string) {
	address := fmt.Sprintf("localhost:%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error al iniciar el servicio de escucha:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Nodo %s escuchando en puerto %s\n", name, port)
	fmt.Println("Local address:", listener.Addr())

	for { // Bucle para aceptar conexiones
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go handleConnection(conn)
	}
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
		log.Println("Error al serializar datos:", err)
		return
	}
	// fmt.Printf("Enviando datos al servidor: %s\n", string(jsonData))
	fmt.Fprintln(conn, string(jsonData))
}

func main() {
	port := os.Getenv("PORT")
	name := os.Getenv("NODE_NAME")
	if port == "" {
		log.Fatal("El puerto no está configurado en la variable de entorno PORT")
	}
	startListening(port, name)
}