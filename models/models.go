package models

// PlayResponse represents the JSON response for a single play
type PlayResponse struct {
    Fruits  []string `json:"fruits"`
    Message string   `json:"message"`
}

// Play10Response represents the JSON response for playing 10 times
type Play10Response struct {
    Spins    []PlayResponse `json:"spins"`
    WinCount int            `json:"win_count"`
}