package utils

import (
	"encoding/csv"
	"fmt"
	"PC4/fc"
	"os"
	"net"
	"time"
)

const(
	port = "9005"
)
// LoadDataset carga el archivo CSV y lo convierte en un DataFrame
func LoadDataset(path string) [][] string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return nil
	}
	
	if len(records) == 0 {
		fmt.Println("Archivo CSV vacío")
		return nil
	}
	return records[1:]
}

// PredictFCC compara los tiempos de entrenamiento de fc y fc_c
func PredictFCC(users map[int]fc.User, targetUser int, k int) {
	start := time.Now()
	fmt.Printf("Predicciones para el usuario %d\n", targetUser)
	recommendationsFCC := fc.RecommendItemsC(users, targetUser, k)
	fmt.Printf("Recomendaciones de fc_c: %v\n", recommendationsFCC)
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución de fc_c: %v\n", elapsed)
}

func ServicioEscuchar() {
	//Modo escucha
	server := fmt.Sprintf("localhost:%s", port) // Puerto 9005
    ln, err := net.Listen("tcp", server)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
        return
    }
    defer ln.Close()
    for { // Modo constante de escucha
        con, err := ln.Accept()
        if err != nil {
            fmt.Println("Error al aceptar la conexión:", err)
            continue
        }
		// Manejar concurrentemente las conexiones
        go fc.Handle(con)
    }
}