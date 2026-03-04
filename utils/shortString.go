package utils

import "fmt"

const SHORT_STRING_SIZE = 10

func GetShortString(input string) string {
	if len(input) > 0 {
		input = input[:SHORT_STRING_SIZE] + ".."
	}
	return fmt.Sprintf("%s", input)
}
