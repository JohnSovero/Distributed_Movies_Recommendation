package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Global WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func getAllMovies(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(movies, "", "  ")
	resp.Write(jsonBytes)
	log.Println("Calling getAllMovies")
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

func getRecommendations(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	log.Println("Calling getRecommendations")
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
	genre := vars["genre"]

	conn, err := net.Dial("tcp", "server:9000")
	if err != nil {
		http.Error(resp, "Error connecting to the server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Create recommendation request with user ID and number of recommendations
	recReq := RecommendationRequest{
		UserID: id,
		NumRec: numRec,
		Genre:  genre,
	}

	requestToServer, err := json.Marshal(recReq)
	if err != nil {
		http.Error(resp, "Error creating request", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(conn, string(requestToServer))

	bf := bufio.NewReader(conn)
	moviesRec, err := bf.ReadString('\n')
	if err != nil {
		http.Error(resp, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Trim the newline and unmarshal the response into a JSON object
	moviesRec = strings.TrimSpace(moviesRec)
	var recommendations []Movie
	err = json.Unmarshal([]byte(moviesRec), &recommendations)
	if err != nil {
		http.Error(resp, "Error parsing recommendations", http.StatusInternalServerError)
		return
	}

	// Send the recommendations back as JSON
	respBytes, err := json.MarshalIndent(recommendations, "", "  ")
	if err != nil {
		http.Error(resp, "Error serializing recommendations", http.StatusInternalServerError)
		return
	}
	resp.Write(respBytes)
}

// Helper function to send the "above average" recommendation request to the server
func sendAboveAverageRequest(conn *websocket.Conn, userID int) {
	// Dial the TCP server
	tcpConn, err := net.Dial("tcp", "server:9000")
	if err != nil {
		log.Println("Error connecting to TCP server:", err)
		return
	}
	defer tcpConn.Close()

	// Prepare request data for "above average" recommendations
	req := RecommendationRequest{
		UserID: userID,
		NumRec: 5, // or any number you want to set for recommendations
		Genre:  "Random",
	}
	reqBytes, _ := json.Marshal(req)
	fmt.Fprintln(tcpConn, string(reqBytes)) // Send request to the server

	// Receive and parse the server's response
	reader := bufio.NewReader(tcpConn)
	respData, _ := reader.ReadString('\n')
	respData = strings.TrimSpace(respData)

	// Parse response and send it over WebSocket
	var recommendations []Movie
	if err := json.Unmarshal([]byte(respData), &recommendations); err == nil {
		conn.WriteJSON(recommendations)
	} else {
		log.Println("Error parsing server response:", err)
	}
}

func wsGetAboveAverageRecommendations(resp http.ResponseWriter, req *http.Request) {
	// Get the user ID from the query parameters
	userIDStr := req.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(resp, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Log or use the userID
	log.Printf("Received User ID: %d", userID)

	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Channel to stop the loop
	stop := make(chan struct{})

	// Start a goroutine to listen for WebSocket close events
	go func() {
		// This will block until the connection is closed or an error occurs
		for {
			_, _, err := conn.NextReader()
			if err != nil {
				close(stop) // Signal to stop the ticker loop
				break
			}
		}
	}()

	// Log or use the userID for further processing
	log.Printf("Assigned User ID: %d", userID)

	// Send initial request
	sendAboveAverageRequest(conn, userID)

	// Set up a ticker to send requests every 1 Minute
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	// Loop to send requests periodically
	for {
		select {
		case <-ticker.C:
			// Send request every 2 minutes
			sendAboveAverageRequest(conn, userID)
		case <-stop:
			// Stop the loop if the WebSocket is closed
			log.Println("WebSocket closed by client")
			return
		}
	}
}
