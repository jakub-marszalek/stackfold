/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"fmt"
	"strings"
)

var UIConfig = struct {
	LineWidth      int
	ShortLineWidth int
	LineChar       string
	EdgeChar       string
	TitleColor     string
	SeparatorColor string
	Reset          string
	Yellow         string
	Green          string
	Red            string
	Cyan           string
}{
	LineWidth:      80,
	ShortLineWidth: 55,
	LineChar:       "-",
	EdgeChar:       "+",
	TitleColor:     "\033[0;36m", // Cyan
	SeparatorColor: "\033[90m",   // Dark Grey
	Reset:          "\033[0m",
	Yellow:         "\033[1;33m",
	Green:          "\033[0;32m",
	Red:            "\033[0;31m",
	Cyan:           "\033[0;36m",
}

const asciiLogo = `    _____ __             __   ____      __    __ 
   / ___// /_____ ______/ /__/ __/___  / /___/ / 
   \__ \/ __/ __ '/ ___/ //_/ /_/ __ \/ / __  /  
  ___/ / /_/ /_/ / /__/ ,< / __/ /_/ / / /_/ /   
 /____/\__/\__,_/\___/_/|_/_/  \____/_/\__,_/`

func DrawLine(color ...string) {
	chosenColor := UIConfig.TitleColor
	if len(color) > 0 {
		chosenColor = color[0]
	}

	innerWidth := UIConfig.LineWidth - 2
	line := UIConfig.EdgeChar + strings.Repeat(UIConfig.LineChar, innerWidth) + UIConfig.EdgeChar
	fmt.Printf("%s%s%s\n", chosenColor, line, UIConfig.Reset)
}

func PrintHeader(text string) {
	DrawLine()

	padding := (UIConfig.LineWidth - len(text)) / 2
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("%s%*s%s%s\n", UIConfig.TitleColor, padding, "", text, UIConfig.Reset)
	DrawLine()
}

func PrintSection(text string) {
	fmt.Printf("%s%s%s\n", UIConfig.Yellow, text, UIConfig.Reset)
}

func LogError(message string) {
	fmt.Printf("%s%s%s\n", UIConfig.Red, message, UIConfig.Reset)
}

func LogWarning(message string) {
	fmt.Printf("%s%s%s\n", UIConfig.Yellow, message, UIConfig.Reset)
}

func LogSuccess(message string) {
	fmt.Printf("%s%s%s\n", UIConfig.Green, message, UIConfig.Reset)
}
