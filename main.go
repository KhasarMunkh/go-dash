package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	base    = "https://api.pandascore.co"
	api_key string
	lrID    = 135916
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

	mux.HandleFunc("/api/teams", TeamSearchHandler)
	mux.HandleFunc("/api/upcoming-matches", GetUpcomingMatchesHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// TeamSearchHandler handles requests to search for teams by name
func TeamSearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for team search")

	q := strings.TrimSpace(r.URL.Query().Get("q")) // e.g., "t1,gen-g"
	game := r.URL.Query().Get("game")              // e.g., "lol" or "csgo"
	idsQS := r.URL.Query().Get("ids")              // e.g., "135916,13917, 2002"
	limit := clamp(parseInt(r.URL.Query().Get("limit"), 10), 1, 50)
	page := clamp(parseInt(r.URL.Query().Get("page"), 1), 1, 100)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Set a timeout for the request, it will return an error if it takes too long
	defer cancel()

	// 1) If client asks for specific team IDs, we will return those teams
	if idsQS != "" {
		ids, err := parseIDs(idsQS)
		if err != nil {
			http.Error(w, "Invalid team IDs", http.StatusBadRequest)
			return
		}
		teams, err := getTeamsByIDs(ctx, game, ids)
		if err != nil {
			http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, teams)
		return
	}
	// 2) If no specific IDs are provided, we can return a default set of teams
	if idsQS == "" {
		// return emtpy or popular suggestions, for now empty
		writeJSONResponse(w, []Team{})
		return
	}
	// 3) Normal case: client asks for teams by game
	if len(q) < 2 {
		http.Error(w, "query too short", http.StatusBadRequest)
		return
	}
	teams, err := getTeams(ctx, q, game, limit, page)
	if err != nil {
		http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, teams)
}

func getTeams(ctx context.Context, q string, game string, limit int, page int) ([]Team, error) {
	path := "/teams"
	if game != "" {
		path = "/" + game + path // e.g., "/lol/teams"
	}

	u := mustURL(base, path) // e.g., "https://api.pandascore.co/lol/teams"

	query := u.Query()
	query.Set("page[size]", strconv.Itoa(limit))
	query.Set("page[number]", strconv.Itoa(page))
	query.Set("filter[name]", q) // e.g., "t1,gen-g"
	u.RawQuery = query.Encode()  // e.g., "https://api.pandascore.co/lol/teams?page[size]=5&page[number]=1&filter[name]=t1,gen-g"

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	defer res.Body.Close()

	var t []Team
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode teams response: %w", err)
	}
	if len(t) == 0 {
		fmt.Println("No teams found for the given query")
		t = []Team{} // Return an empty slice if no teams are found
	}
	return t, nil
}

func getTeamsByIDs(ctx context.Context, game string, ids []int) ([]Team, error) {
	path := "/teams"
	if game != "" {
		path = "/" + game + path // e.g., "/lol/teams"
	}

	u := mustURL(base, path) // e.g., "https://api.pandascore.co/lol/teams"
	query := u.Query()
	query.Set("filter[id]", joinInts(ids, ",")) // e.g., "135916,13917,2002"
	u.RawQuery = query.Encode() // e.g., "https://api.pandascore.co/lol/teams?filter[id]=135916,13917,2002"

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams by IDs: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		//b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to fetch teams, status code: %d", res.StatusCode)
	}

	var t []Team
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode teams response: %w", err)
	}

	return t, nil
}

func GetUpcomingMatchesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Set a timeout for the request
	defer cancel()


	game := strings.TrimSpace(r.URL.Query().Get("game")) // e.g., "lol"
	idsQS := strings.TrimSpace(r.URL.Query().Get("ids")) // e.g., "135916,13917,2002"
	var teamIDs []int
	if idsQS != "" {
		var err error 
		teamIDs, err = parseIDs(idsQS)
		if err != nil {
			http.Error(w, "Invalid team IDs", http.StatusBadRequest)
			return
		}
	}

	m, err := RequestMatches(ctx, game, teamIDs) 
	if err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, m)
}

func RequestMatches(ctx context.Context, game string, teamIDs []int) ([]Match, error) {
	fmt.Println("Requesting matches...")
	path := "/matches/upcoming"
	if game != "" {
		path = "/" + game + path // e.g., "/lol/matches/upcoming"
	}

	u := mustURL(base, path)
	query := u.Query()
	query.Set("page[size]", "100") // Set the page size to 100
	query.Set("page[number]", "1") // Set the page number to 1
	query.Set("sort", "scheduled_at") // Sort by scheduled date

	if len(teamIDs) > 0 {
		query.Set("filter[opponent_id]", joinInts(teamIDs, ",")) // e.g., "135916"
	}

	u.RawQuery = query.Encode() // e.g., "https://api.pandascore.co/lol/matches/upcoming?page[size]=100&page[number]=1&sort=scheduled_at&filter[opponent_id]=135916"
	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to fetch matches, status code: %d", res.StatusCode, string(b))
	}

	var m []Match

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}
	if m == nil {
		m = []Match{}
	}

	return m, nil
}

// func RequestTournament(game string) {
//	path := "/lol/series/running"
//	u := mustURL(base, path)
//	req, _ := http.NewRequest("GET", u.String(), nil)
//
//	req.Header.Add("accept", "application/json")
//	req.Header.Add("authorization", "Bearer "+api_key)
//
//	res, _ := http.DefaultClient.Do(req)
//
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//
//	fmt.Println(string(body))
// }

// func RequestUpcomingByTeam(teamID int) ([]Match, error) {
//	fmt.Printf("Requesting upcoming matches for team ID: %d\n", teamID)
//	u := fmt.Sprintf("%s/lol/matches/upcoming?filter[opponent_id]=%d", base, teamID)
//
//	req, _ := http.NewRequest("GET", u, nil)
//
//	req.Header.Add("accept", "application/json")
//	req.Header.Add("authorization", "Bearer "+api_key)
//
//	res, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	defer res.Body.Close()
//
//	var m []Match
//
//	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
//		return nil, err
//	}
//
//	return m, nil
// }

// Helpers
func writeJSONResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func parseIDs(idsQS string) ([]int, error) {
	parts := strings.Split(idsQS, ",")
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p == "" {
			continue
		} // Skip empty parts
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	return ids, nil
}

func mustURL(base, path string) *url.URL {
    b := strings.TrimRight(base, "/")
    p := "/" + strings.TrimLeft(path, "/")
    u, err := url.Parse(b + p)
    if err != nil { panic(err) }
    return u
}

func joinInts(xs []int, sep string) string {
    if len(xs) == 0 { return "" }
    var b strings.Builder
    for i, n := range xs {
        if i > 0 { b.WriteString(sep) }
        b.WriteString(strconv.Itoa(n))
    }
    return b.String()
}

func parseInt(s string, def int) int {
    s = strings.TrimSpace(s)
    if s == "" { return def }
    n, err := strconv.Atoi(s)
    if err != nil { return def }
    return n
}

func clamp(n, lo, hi int) int {
    if n < lo { return lo }
    if n > hi { return hi }
    return n
}
