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

	"github.com/gorilla/mux"
)

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

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		http.Error(resp, "Error connecting to the server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Create recommendation request with user ID and number of recommendations
	recReq := RecommendationRequest{
		UserID: id,
		NumRec: numRec,
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
	var recommendations []int
	err = json.Unmarshal([]byte(moviesRec), &recommendations)
	if err != nil {
		http.Error(resp, "Error parsing recommendations", http.StatusInternalServerError)
		return
	}

	// Send the recommendations back as JSON
	respBytes, err := json.Marshal(recommendations)
	if err != nil {
		http.Error(resp, "Error serializing recommendations", http.StatusInternalServerError)
		return
	}
	resp.Write(respBytes)
}
