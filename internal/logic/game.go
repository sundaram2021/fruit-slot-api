package logic

import (
    "crypto/rand"
    "fmt"
    "math/big"
)

const (
    WinMessage  = "You win!"
    LoseMessage = "Try again!"
)

var FruitsPool = []string{"Cherry", "Lemon", "Orange", "Grape", "Watermelon", "Pineapple"}

// getRandomFruits selects 'count' number of fruits randomly from the pool.
// It allows for the same fruit to be selected multiple times.
func GetRandomFruits(count int) ([]string, error) {
    if count <= 0 {
        return []string{}, nil
    }

    selectedFruits := make([]string, count)
    poolSize := int64(len(FruitsPool))

    for i := 0; i < count; i++ {
        // Generate a secure random index
        indexBig, err := rand.Int(rand.Reader, big.NewInt(poolSize))
        if err != nil {
            return nil, fmt.Errorf("failed to generate random number: %w", err)
        }
        index := indexBig.Int64()
        selectedFruits[i] = FruitsPool[index]
    }
    return selectedFruits, nil
}

// checkWin determines if the player wins based on the fruits.
// A win occurs if 2 or more fruits are the same.
func CheckWin(fruits []string) bool {
    if len(fruits) < 2 {
        return false // can't have less than 2 fruits, wining criteria
    }

    // Count occurrences of each fruit
	// map is best Data Structure for this , complexity O(1)
    fruitCounts := make(map[string]int)
    for _, fruit := range fruits {
        fruitCounts[fruit]++
        if fruitCounts[fruit] >= 2 { // no need to range further , we got our ans
            return true
        }
    }
    return false // No fruit appeared 2 or more times
}