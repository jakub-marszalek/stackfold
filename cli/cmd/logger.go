package cmd

import (
	"fmt"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
)

func LogError(message string) {
	fmt.Printf("%s%s%s\n", ColorRed, message, ColorReset)
}

func LogWarning(message string) {
	fmt.Printf("%s%s%s\n", ColorYellow, message, ColorReset)
}
