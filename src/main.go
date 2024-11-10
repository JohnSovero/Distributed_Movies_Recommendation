package main

import (
	"PC4/fc"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func serverStartListening(port string) {
	address := "localhost:" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error al iniciar el servicio de escucha:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go serverHandleConnection(conn)
	}
}

func serverHandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			return
		}
		fmt.Println("Mensaje recibido:", message)
	}
}

func main() {
	// Leer archivo de recomendación de películas
	pathRatings := "dataset/ratings25.csv"
	pathMovies := "dataset/movies25.csv"
	fmt.Println("\nLeyendo archivos de datos...")
	fmt.Println("--------------------------------")
	fmt.Println("Detalle de la información procesada:")
	ratings, err := fc.ReadRatingsFromCSV(pathRatings)
	if err != nil {
		log.Fatalf("Error leyendo los ratings del csv: %v", err)
	}
	movies, err := fc.ReadMoviesFromCSV(pathMovies)
	if err != nil {
		log.Fatalf("Error leyendo las películas del csv: %v", err)
	}

	var userId int = -1
	var numRecommendations int = 5 // Valor por defecto
	reader := bufio.NewReader(os.Stdin)

	serverPort := "9015"
	startListening(serverPort)

	// for {
	// 	fmt.Println("--------------------------------")
	// 	fmt.Println("||||||||||||||||||||||||||||||||")
	// 	fmt.Println("--------------------------------")
	// 	if userId == -1 {
	// 		fmt.Println("\nUsuario: No especificado")
	// 	} else {
	// 		fmt.Println("\nUsuario:", userId)
	// 	}
	// 	fmt.Println("Número de recomendaciones:", numRecommendations)
	// 	fmt.Println("\nMenú de opciones: ------------")
	// 	fmt.Println()
	// 	fmt.Println("\t1. Ingresar ID del usuario")
	// 	fmt.Println("\t2. Indicar cuántas películas recomendar")
	// 	fmt.Println("\t3. Predecir recomendaciones")
	// 	fmt.Println("\t4. Salir")
	// 	fmt.Println()
	// 	fmt.Println("--------------------------------")
	// 	fmt.Print("Seleccione una opción: ")
	// 	optionStr, _ := reader.ReadString('\n')
	// 	fmt.Println("--------------------------------")

	// 	optionStr = strings.TrimSpace(optionStr)
	// 	option, err := strconv.Atoi(optionStr)
	// 	if err != nil {
	// 		fmt.Println("Opción inválida. Intente de nuevo.")
	// 		continue
	// 	}

	// 	switch option {
	// 	case 1:
	// 		fmt.Print("Ingrese el ID del usuario: ")
	// 		userIdStr, _ := reader.ReadString('\n')
	// 		userIdStr = strings.TrimSpace(userIdStr)
	// 		userId, err = strconv.Atoi(userIdStr)
	// 		if err != nil {
	// 			fmt.Println("ID de usuario inválido. Intente de nuevo.")
	// 			userId = -1
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	case 2:
	// 		fmt.Print("Ingrese el número de películas a recomendar: ")
	// 		numRecsStr, _ := reader.ReadString('\n')
	// 		numRecsStr = strings.TrimSpace(numRecsStr)
	// 		numRecommendations, err = strconv.Atoi(numRecsStr)
	// 		if err != nil || numRecommendations <= 0 {
	// 			fmt.Println("Número de recomendaciones inválido. Usando valor por defecto (5).")
	// 			numRecommendations = 5
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	case 3:
	// 		if userId == -1 {
	// 			fmt.Println("Primero debe ingresar un ID de usuario válido.")
	// 			time.Sleep(1 * time.Second)
	// 		} else {
	// 			fc.PredictFC(ratings, userId, numRecommendations, movies)
	// 		}
	// 	case 4:
	// 		fmt.Println("Saliendo...")
	// 		return
	// 	default:
	// 		fmt.Println("Opción inválida. Intente de nuevo.")
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }
}
