package utils

import (
	"encoding/csv"
	"fmt"
	"PC4/fc"
	"os"
	"time"
	"net"
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

// MeasureExecutionTime mide el tiempo de ejecución de una función
func MeasureExecutionTime(name string, f func()) {
    start := time.Now()
    f()
    duration := time.Since(start)
    fmt.Printf("Tiempo de ejecución para %s: %v\n", name, duration)
}

// PredictFCC compara los tiempos de entrenamiento de fc y fc_c
func PredictFCC(users map[int]fc.User, targetUser int, k int, ln net.Listener) {
	fmt.Printf("Predicciones para el usuario %d\n", targetUser)
	recommendationsFCC := fc.RecommendItemsC(users, targetUser, k, ln)
	fmt.Printf("Recomendaciones de fc_c: %v\n", recommendationsFCC)
}
