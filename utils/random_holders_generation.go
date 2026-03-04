package utils

import (
	"math/rand"
)

func GenerateRandomHoldersUniform(min int, max int, size int) []int {
	sampledNumbers := make(map[int]bool)
	result := make([]int, 0, size)
	for len(result) < size {
		num := rand.Intn(max-1) + min // Generate a random number in [min, max]
		//zap.S().Infoln(num)
		if !sampledNumbers[num] {
			sampledNumbers[num] = true
			result = append(result, num)
		}
	}
	return result
}
