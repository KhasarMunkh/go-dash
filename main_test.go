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

/*
func TestRequestMatches(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    []Match
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := RequestMatches()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RequestMatches() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RequestMatches() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("RequestMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

*/
