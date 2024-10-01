package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var inputReader = bufio.NewReader(os.Stdin)
var dict = []string{
	"Zombie",
	"Gopher",
	"United States of America",
	"Indonesia",
	"Nazism",
	"Apple",
	"Samsung",
	"Programming",
}

func main() {
	targetWord := getRandomWord()
	guessedLetters := initializeGeussedWords(targetWord)
	hangmanState := 0
	fmt.Println("Random word:", targetWord)

	for !isGameOver(targetWord, guessedLetters, hangmanState) {
		printGameState(targetWord, guessedLetters, hangmanState)
		input := readInput()
		if len(input) != 1 {
			fmt.Println("Invalid input. Please use letters only....")
			continue
		}

		letter := rune(input[0])
		if isCorrectGuess(targetWord, letter) {
			guessedLetters[letter] = true
		} else {
			hangmanState++
		}
	}

	fmt.Print("Game Over!")
	if isWordGuessed(targetWord, guessedLetters) {
		fmt.Println("You win!")
	} else if isHangmanComplete(hangmanState) {
		printGameState(targetWord, guessedLetters, hangmanState)
		fmt.Println("You lose!")
	} else {
		panic("invalid state. Game over there is no winner")
	}
}

func getRandomWord() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	return dict[rng.Intn(len(dict))]
}

func initializeGeussedWords(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	guessedLetters[unicode.ToLower(rune(targetWord[0]))] = true
	guessedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true

	return guessedLetters
}

func isGameOver(
	targetWord string,
	guesssedLetters map[rune]bool,
	hangmanState int,
) bool {
	return isWordGuessed(targetWord, guesssedLetters) || isHangmanComplete(hangmanState)
}

func isWordGuessed(targetWord string, guessedLetters map[rune]bool) bool {
	for _, ch := range targetWord {
		if !guessedLetters[unicode.ToLower(ch)] {
			return false
		}
	}
	return true
}

func isHangmanComplete(hangmanState int) bool {
	return hangmanState >= 9
}

func printGameState(targetWord string, guessedLetters map[rune]bool, hangmanState int) {
	fmt.Println(getWordGuessingProgress(targetWord, guessedLetters))
	fmt.Println()
	fmt.Println(randerHangman(hangmanState))
}

func getWordGuessingProgress(targetWord string, guessedLetters map[rune]bool) string {
	result := ""
	for _, ch := range targetWord {
		if ch == ' ' {
			result += " "
		} else if guessedLetters[unicode.ToLower(ch)] {
			result += fmt.Sprintf("%c", ch)
		} else {
			result += "_"
		}
		result += " "
	}
	return result
}

func randerHangman(hangmanState int) string {
	data, err := os.ReadFile(fmt.Sprintf("states/hangman%d", hangmanState))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func readInput() string {
	fmt.Print("> ")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(input)
}

func isCorrectGuess(targetWord string, letter rune) bool {
	return strings.ContainsRune(targetWord, letter)
}
