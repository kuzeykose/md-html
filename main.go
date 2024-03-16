package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Open the file
	file, err := os.Open("test.md")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	for scanner.Scan() {
		// Print each line to the console
		line := scanner.Text()
		tokenizer(line)
		fmt.Println()
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

type Token struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func tokenizer(line string) {
	tokens := []Token{}

	runes := []rune(line)

	i := 0
	for i < len(runes) {

		// if i == 0 && runes[i+1] == ' ' && runes[i] == '#' {
		// tokens = append(tokens, Token{Type: "identifier", Value: "#"})
		// }
	}

	for _, elem := range tokens {
		fmt.Println("Type:", elem.Type, "Value:", elem.Value)
	}

}

func t(runes []rune, i int) {

}

// func isIdentifier(char rune, pos int) Token{
// 	switch char {
// 	case '#':
// 		isHeading()
// 	}
// }
