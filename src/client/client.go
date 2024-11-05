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
)

const(
    portServer = "9005"
)

var server string
var portClient string
type User struct {
	ID      int
	Ratings map[int]float64
}

// Estructura para enviar datos al servidor
type ToServer struct {
    Similarity float64 `json:"similarity"`
    UserID     string  `json:"userID"`
}

// Estructura para recibir datos del cliente
type ClientData struct {
    User1 map[int]float64 `json:"user1"`
    User2 map[int]User `json:"user2"`
}

// Maneja la conexión del cliente
func handleClient(conn net.Conn) {
    str, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    str = strings.TrimSpace(str)

    // Deserializar JSON a ClientData
    var data ClientData
    json.Unmarshal([]byte(str), &data)
    
    // Calcular similitud coseno entre los usuarios y enviar al servidor
    for _, user := range data.User2 {
        go func(user User) {
            result := cosineSimilarity(data.User1, user.Ratings)
            sendToServer(result, strconv.Itoa(user.ID), conn)
        }(user)
    }
}

func servicioEscuchar(port string){
    // Configurar el cliente
    dirClient := fmt.Sprintf("localhost:%s", port)
    ln, err := net.Listen("tcp", dirClient)
    if err != nil {
        fmt.Println("Error al iniciar el cliente:", err)
        return
    }
    defer ln.Close()
    for { // Modo constante de escucha
        con, err := ln.Accept()
        if err != nil {
            fmt.Println("Error al aceptar la conexión:", err)
            continue
        }
        go handleClient(con)
    }
}

// Lee la entrada del usuario
func getUserInput() string {
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

// Función para calcular la similitud coseno entre dos usuarios
func cosineSimilarity(user1 map[int]float64, user2 map[int]float64) float64 {
    dotProduct := 0.0
    normA := 0.0
    normB := 0.0

    for itemID, rating1 := range user1 {
        if rating2, exists := user2[itemID]; exists {
            dotProduct += rating1 * rating2
            normA += rating1 * rating1
            normB += rating2 * rating2
        }
    }
    result := 0.0
    // Evitar división por cero
    if normA == 0 || normB == 0 {
        result = 0.0
    } else {
        result = dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
    }
    return result
}

// Envía resultados al servidor
func sendToServer(similarity float64, userID string, conn net.Conn) {
    defer conn.Close()
    message := ToServer{
        Similarity: similarity,
        UserID:     userID,
    }
    // serializar
    jsonData, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshaling to JSON:", err)
        return
    }

    fmt.Printf("Sending JSON: %s\n", jsonData)
    fmt.Fprintln(conn, string(jsonData))
}

func main() {
    server = fmt.Sprintf("localhost:%s", portServer)
    fmt.Print("Enter port: ")
    portClient := getUserInput()
    servicioEscuchar(portClient)
}