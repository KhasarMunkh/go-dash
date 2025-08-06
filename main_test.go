package main

import (
	"testing"
)

func TestRequestMatches(t *testing.T) {
	matches, err := RequestMatches()
	if err != nil {
		t.Fatalf("Failed to fetch matches: %v", err)
	}

	if len(matches) == 0 {
		t.Errorf("Expected at least one match, got %d", len(matches))
	}

	for _, match := range matches {
		t.Logf("Match ID: %d, Date: %s, League: %s", match.Id, match.Date, match.League.Name)
		for _, opponent := range match.Opponents {
			t.Logf(" - %s", opponent.Opponent.Name)
		}
	}
}
