package handlers

import (
	"github.com/sundaram2021/fruit-slot-api/internal/logic"
	"github.com/sundaram2021/fruit-slot-api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Play handles the GET /play endpoint.
// It generates 3 random fruits, checks for a win, and returns the result.
func Play(c *gin.Context) {
	log.Println("Received request for /play")

	// GetRandomFruits: will get  3 random fruits
	fruits, err := logic.GetRandomFruits(3)
	if err != nil {
		log.Printf("Error generating fruits: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	isWin := logic.CheckWin(fruits)
	message := logic.LoseMessage
	if isWin {
		message = logic.WinMessage
	}

	response := models.PlayResponse{
		Fruits:  fruits,
		Message: message,
	}

	log.Printf("Spin result: %+v, Win: %t", fruits, isWin)
	c.JSON(http.StatusOK, response)
}

// Play10 handles the GET /play/10 endpoint.
// It generates 10 sets of 3 random fruits, checks wins for each,
// and returns the results along with the total win count.
func Play10(c *gin.Context) {
	log.Println("Received request for /play/10")

	const numSpins = 10
	spins := make([]models.PlayResponse, numSpins)
	winCount := 0

	for i := 0; i < numSpins; i++ {
		fruits, err := logic.GetRandomFruits(3)
		if err != nil {
			log.Printf("Error generating fruits for spin %d: %v", i+1, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return 
		}

		isWin := logic.CheckWin(fruits)
		message := logic.LoseMessage
		if isWin {
			message = logic.WinMessage
			winCount++
		}

		spins[i] = models.PlayResponse{
			Fruits:  fruits,
			Message: message,
		}
		// Log each spin if needed, or aggregate logs
		log.Printf("Spin %d result: %+v, Win: %t", i+1, fruits, isWin)
	}

	response := models.Play10Response{
		Spins:    spins,
		WinCount: winCount,
	}

	log.Printf("Completed 10 spins. Total wins: %d", winCount)
	c.JSON(http.StatusOK, response)
}
