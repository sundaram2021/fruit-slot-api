package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/sundaram2021/fruit-slot-api/internal/logic"
	"github.com/sundaram2021/fruit-slot-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	// Switch to test mode so Gin doesn't print logs
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/play", Play)
	r.GET("/play/10", Play10)
	return r
}

func TestPlayEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/play", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PlayResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Basic structure checks
	assert.Len(t, response.Fruits, 3)
	assert.Contains(t, []string{logic.WinMessage, logic.LoseMessage}, response.Message)

	// Validate fruits are from the pool
	for _, fruit := range response.Fruits {
		assert.Contains(t, logic.FruitsPool, fruit)
	}
}

func TestPlay10Endpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/play/10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Play10Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Basic structure checks
	assert.Len(t, response.Spins, 10)
	assert.GreaterOrEqual(t, response.WinCount, 0)
	assert.LessOrEqual(t, response.WinCount, 10)

	// Validate each spin
	for _, spin := range response.Spins {
		assert.Len(t, spin.Fruits, 3)
		assert.Contains(t, []string{logic.WinMessage, logic.LoseMessage}, spin.Message)
		for _, fruit := range spin.Fruits {
			assert.Contains(t, logic.FruitsPool, fruit)
		}
	}
}

// Test the core logic directly for more control
func TestCheckWinLogic(t *testing.T) {
	tests := []struct {
		name     string
		fruits   []string
		expected bool
	}{
		{"All Different", []string{"Cherry", "Lemon", "Orange"}, false},
		{"Two Same", []string{"Grape", "Grape", "Watermelon"}, true},
		{"Three Same", []string{"Pineapple", "Pineapple", "Pineapple"}, true},
		{"Empty", []string{}, false},
		{"One Fruit", []string{"Cherry"}, false},
		{"Two Different", []string{"Lemon", "Orange"}, false},
		{"Two Same at End", []string{"Cherry", "Lemon", "Lemon"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logic.CheckWin(tt.fruits)
			assert.Equal(t, tt.expected, result, "checkWin(%v) = %v; expected %v", tt.fruits, result, tt.expected)
		})
	}
}

// written concurrent test for play endpoint
func TestConcurrentPlayEndpoint(t *testing.T) {
	router := setupRouter()
	const concurrentRequests = 100

	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/play", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf("expected status 200, got %d", w.Code)
				return
			}

			var response models.PlayResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				errors <- fmt.Errorf("failed to unmarshal response: %v", err)
				return
			}

			if len(response.Fruits) != 3 {
				errors <- fmt.Errorf("expected 3 fruits, got %d", len(response.Fruits))
				return
			}

			if response.Message != logic.WinMessage && response.Message != logic.LoseMessage {
				errors <- fmt.Errorf("unexpected message: %s", response.Message)
				return
			}

			for _, fruit := range response.Fruits {
				if !contains(logic.FruitsPool, fruit) {
					errors <- fmt.Errorf("unexpected fruit: %s", fruit)
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// helper function
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// written concurrent test for play/10 endpoint
func TestConcurrentPlay10Endpoint(t *testing.T) {
	router := setupRouter()
	const concurrentRequests = 100

	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/play/10", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf("expected status 200, got %d", w.Code)
				return
			}

			var response models.Play10Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				errors <- fmt.Errorf("failed to unmarshal response: %v", err)
				return
			}

			if len(response.Spins) != 10 {
				errors <- fmt.Errorf("expected 10 spins, got %d", len(response.Spins))
				return
			}

			if response.WinCount < 0 || response.WinCount > 10 {
				errors <- fmt.Errorf("win count out of bounds: %d", response.WinCount)
				return
			}

			for _, spin := range response.Spins {
				if len(spin.Fruits) != 3 {
					errors <- fmt.Errorf("expected 3 fruits in spin, got %d", len(spin.Fruits))
					return
				}

				if spin.Message != logic.WinMessage && spin.Message != logic.LoseMessage {
					errors <- fmt.Errorf("unexpected message: %s", spin.Message)
					return
				}

				for _, fruit := range spin.Fruits {
					if !contains(logic.FruitsPool, fruit) {
						errors <- fmt.Errorf("unexpected fruit: %s", fruit)
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

