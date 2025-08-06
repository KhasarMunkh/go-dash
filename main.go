package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	base    = "https://api.pandascore.co"
	api_key string
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
}

func RequestMatches() ([]Match, error) {
	fmt.Println("Requesting matches...")
	url := base + "/lol/matches/upcoming?filter[opponent_id]=t1,gen-g"

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

type League struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type Match struct {
	Name      string          `json:"name"`
	Id        int             `json:"id"`
	Date      string          `json:"scheduled_at"`
	BeginAt   time.Time       `json:"begin_at"`
	Opponents []OpponentEntry `json:"opponents"`
	League    League          `json:"league"`
}

type OpponentEntry struct {
	Opponent Team `json:"opponent"`
}

type Team struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Acronym  string `json:"acronym"`
	ImageURL string `json:"image_url"`
}
