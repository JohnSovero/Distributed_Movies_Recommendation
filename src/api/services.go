package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func getAllMovies(resp http.ResponseWriter, req *http.Request) {
	log.Println("Calling getAllMovies")
	resp.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		http.Error(resp, "Error serializing movies", http.StatusInternalServerError)
		return
	}
	resp.Write(jsonBytes)
	log.Println("getAllMovies called")
}

func getAllUsers(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(users, "", "  ")
	resp.Write(jsonBytes)
	log.Println("Calling getAllUsers")
}

func getMovieByID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	id := vars["id"]

	var wg sync.WaitGroup
	var mu sync.Mutex
	found := false

	for _, movie := range movies {
		wg.Add(1)
		go func(movie Movie) {
			defer wg.Done()
			if strconv.Itoa(movie.MovieID) == id {
				mu.Lock()
				if !found {
					found = true
					jsonBytes, _ := json.MarshalIndent(movie, "", "  ")
					resp.Write(jsonBytes)
				}
				mu.Unlock()
			}
		}(movie)
	}

	wg.Wait()
	if !found {
		http.Error(resp, "Movie not found", http.StatusNotFound)
	}
}

func getRecommendationsWS(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(resp, "Invalid user ID", http.StatusBadRequest)
		return
	}
	numRec, err := strconv.Atoi(vars["numRec"])
	if err != nil {
		http.Error(resp, "Invalid number of recommendations", http.StatusBadRequest)
		return
	}

	// Updating a WebSocket
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	// Connecting to the server
	tcpConn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Println("Error connecting to recommendation server:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error connecting to recommendation server"))
		return
	}
	defer tcpConn.Close()

	recReq := RecommendationRequest{UserID: id, NumRec: numRec}
	requestToServer, err := json.Marshal(recReq)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error creating request"))
		return
	}

	// Send the request to the server
	fmt.Fprintln(tcpConn, string(requestToServer))

	// Reading the response from the server
	bf := bufio.NewReader(tcpConn)
	moviesRec, err := bf.ReadString('\n')
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error reading response"))
		return
	}

	// Send the response to the client
	conn.WriteMessage(websocket.TextMessage, []byte(moviesRec))
}
