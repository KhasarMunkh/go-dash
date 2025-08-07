package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	base    = "https://api.pandascore.co"
	api_key string
	lrID = 135916
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key = os.Getenv("PANDA_KEY")
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs) // Serve static files from the public directory when accessing the root URL

	mux.HandleFunc("/api/upcoming-matches", GetUpcomingMatchesHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func GetUpcomingMatchesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for upcoming matches")
	m, err := RequestMatches()
	if err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, "Failed to encode matches", http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully fetched and returned matches")
}

func RequestMatches() ([]Match, error) {
	fmt.Println("Requesting matches...")
	// url := base + "/lol/matches/upcoming?filter[opponent_id]=t1,gen-g"
	 url := base + "/lol/matches/upcoming?page[size]=5&page[number]=1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var m []Match

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func RequestSeries() () {
	fmt.Println("Requesting series...")
	// url := base + "/lol/series/upcoming?page[size]=5&page[number]=1"
	url := base+ "/lol/series/running?sort=begin_at&page=1&per_page=1"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer 2MO8I9xEUefnWHxXsVmZyRxmMtRXlIxyX1ZfyZ1W8EHMnJko6a0")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func RequestUpcomingByTeam(teamID int) ([]Match, error) {
	fmt.Printf("Requesting upcoming matches for team ID: %d\n", teamID)
	url := fmt.Sprintf("%s/lol/matches/upcoming?filter[opponent_id]=%d", base, teamID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var m []Match

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}
