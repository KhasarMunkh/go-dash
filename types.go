package main

import (
	"time"
)

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
