package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadData() {
	loadMoviesFromCSV()
	loadUsersFromCSV()
}

func loadMoviesFromCSV() {
	file, err := os.Open("movies_complete.csv")
	if err != nil {
		log.Fatal("Error while opening movies file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header row
	if _, err := reader.Read(); err != nil {
		log.Fatal("Error while reading header row:", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error while reading movies file:", err)
		}
		movieID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal("Error while parsing movie ID:", err)
		}
		year, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatal("Error while parsing movie year:", err)
		}
		movie := Movie{
			MovieID:    movieID,
			Title:      record[1],
			Genres:     strings.Split(record[2], "|"),
			IMDBLink:   record[3],
			TMDBLink:   record[4],
			Year:       year,
			Overview:   record[7],
			VoteAvg:    record[8],
			PosterPath: record[9],
		}
		movies = append(movies, movie)
	}
}

func loadUsersFromCSV() {
	file, err := os.Open("users.csv")
	if err != nil {
		log.Fatal("Error while opening users file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header row
	if _, err := reader.Read(); err != nil {
		log.Fatal("Error while reading header row:", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error while reading users file:", err)
		}
		userID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal("Error while parsing user ID:", err)
		}
		users = append(users, userID)
	}
}