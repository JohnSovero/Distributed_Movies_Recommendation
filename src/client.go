package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "math"
    "net"
    "os"
    "strings"
)

const(
    port = "9005"
)
var server string

// Estructura para enviar datos al servidor
type ToServer struct {
    Similarity float64 `json:"similarity"`
    UserID     string  `json:"userID"`
}

// Estructura para recibir datos del cliente
type ClientData struct {
    User1 map[int]float64 `json:"user1"`
    User2 map[int]float64 `json:"user2"`
    ID    string          `json:"id"`
}

func main() {
    server = "localhost:9005"
    fmt.Print("Enter port: ")
    port := getUserInput()

    servicioEscuchar(port)
}
// Maneja la conexión del cliente
func handleClient(conn net.Conn) {
    defer conn.Close()
    str, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    str = strings.TrimSpace(str)

    // Deserializar JSON a ClientData
    var data ClientData
    json.Unmarshal([]byte(str), &data)

    // Calcular similitud de coseno
    similarity := cosineSimilarity(data.User1, data.User2)

    // Enviar resultado de vuelta al servidor
    sendToServer(similarity, data.ID)
}

func servicioEscuchar(port string){
    // Configurar el cliente
    localDir := fmt.Sprintf("localhost:%s", port)
    ln, err := net.Listen("tcp", localDir)
    if err != nil {
        fmt.Println("Error al iniciar el cliente:", err)
        return
    }
    defer ln.Close()

    // Aceptar conexiones entrantes
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error al aceptar la conexión:", err)
            continue
        }
        go handleClient(conn)
    }
}
// Lee la entrada del usuario
func getUserInput() string {
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

// Función para calcular la similitud coseno entre dos usuarios
func cosineSimilarity(user1, user2 map[int]float64) float64 {
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

    // Evitar división por cero
    if normA == 0 || normB == 0 {
        return 0
    }
    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Envía resultados al servidor
func sendToServer(similarity float64, userID string) {
    for {
        conn, err := net.Dial("tcp", server)
        if err == nil {
            defer conn.Close()

            message := ToServer{
                Similarity: similarity,
                UserID:     userID,
            }

            jsonData, err := json.Marshal(message)
            if err != nil {
                fmt.Println("Error marshaling to JSON:", err)
                return
            }

            fmt.Printf("Sending JSON: %s\n", jsonData)
            fmt.Fprintln(conn, string(jsonData))
            return
        }
        fmt.Printf("Error connecting to server: %v. Retrying...\n", err)
    }
}
//probar reintentando enviar al server