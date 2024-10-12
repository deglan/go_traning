package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

const SnakeSymbol = 0x2588
const AppleSymbol = 0x25CF
const GameFrameWidth = 30
const GameFrameHight = 20
const GameFrameSymbol = 0x2588

type Point struct {
	row, col int
}

type Snake struct {
	parts          []*Point
	velRow, velCol int
	symbol         rune
}

type Apple struct {
	point  *Point
	symbol rune
}

var screenWidth, screenHeigh int
var screen tcell.Screen
var snake *Snake
var apple *Apple
var score int
var pointsToClear []*Point
var isGameOver bool
var isGamePaused bool
var debugLog string

func main() {

	initScreen()
	screenWidth, screenHeigh = screen.Size()
	initGameState()
	inputChan := initUserInput()

	for !isGameOver {
		handleUserInput(readInpout(inputChan))
		updateState()
		drawState()
		time.Sleep(150 * time.Millisecond)
	}

	printStringCentred(screenHeigh/2, screenWidth/2, "Game Over")
	printStringCentred(screenHeigh/2+1, screenWidth/2, fmt.Sprintf("Your score is %d", score))
	printStringCentred(screenHeigh/2+2, screenWidth/2, "Pres enter to quit")
	screen.Show()
	waitForEnter()
	screen.Fini()
}

func updateState() {
	if isGamePaused {
		return
	}

	updateSnake()
	head := getSnakeHead()
	if head.row == apple.point.row && head.col == apple.point.col {
		updateApple()
	}
}

func updateSnake() {
	head := getSnakeHead()
	snake.parts = append(snake.parts, &Point{
		row: head.row + snake.velRow,
		col: head.col + snake.velCol,
	})

	if !appleInsideSnake() {
		snake.parts = snake.parts[1:]
	} else {
		score++
	}

	if isSnakeHittingWall() || isSnakeEatingItself() {
		isGameOver = true
	}
}

func isSnakeEatingItself() bool {
	head := getSnakeHead()
	for _, p := range snake.parts[:getSnakeHeadIndex()] {
		if p.row == head.row && p.col == head.col {
			return true
		}
	}
	return false
}

func getSnakeHeadIndex() int {
	return len(snake.parts) - 1
}

func isSnakeHittingWall() bool {
	head := getSnakeHead()
	return head.row < 0 ||
		head.row >= GameFrameHight ||
		head.col < 0 ||
		head.col >= GameFrameHight
}

func getSnakeHead() *Point {
	return snake.parts[len(snake.parts)-1]
}

func updateApple() {
	for appleInsideSnake() {
		apple.point.row, apple.point.col = rand.Intn(GameFrameHight), rand.Intn(GameFrameWidth)
	}
}

func appleInsideSnake() bool {
	for _, p := range snake.parts {
		if p.row == apple.point.row && p.col == apple.point.col {
			return true
		}
	}
	return false
}

func drawState() {
	debugLog = fmt.Sprintf("Score: %d", score)
	if isGamePaused {
		return
	}
	clearScreen()
	printString(0, 0, debugLog)
	printGameFrame()
	printSnake()
	printApple()
	screen.Show()
}

func clearScreen() {
	for _, p := range pointsToClear {
		printFilledRectInGameFrame(p.row, p.col, 1, 1, ' ')
	}
	pointsToClear = []*Point{}
}

func initGameState() {
	snake = &Snake{
		parts: []*Point{
			{row: 9, col: 3}, //tail
			{row: 8, col: 3},
			{row: 7, col: 3},
			{row: 6, col: 3},
			{row: 5, col: 3}, //head
		},
		velRow: -1,
		velCol: 0,
		symbol: SnakeSymbol,
	}
	apple = &Apple{
		point:  &Point{row: 10, col: 10},
		symbol: AppleSymbol,
	}
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if key == "Rune[w]" && snake.velRow != 1 {
		snake.velRow = -1
		snake.velCol = 0
	} else if key == "Rune[a]" && snake.velCol != 1 {
		snake.velRow = 0
		snake.velCol = -1
	} else if key == "Rune[s]" && snake.velRow != -1 {
		snake.velRow = 1
		snake.velCol = 0
	} else if key == "Rune[d]" && snake.velCol != -1 {
		snake.velRow = 0
		snake.velCol = 1
	}
}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func readInpout(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func printStringCentred(row, col int, str string) {
	col = col - len(str)/2
	printString(row, col, str)
}

func printString(row, col int, str string) {
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		for _, c := range line {
			screen.SetContent(col, row, c, nil, tcell.StyleDefault)
			col++
		}
		row++
		col = (screenWidth / 2) - (len(line) / 2)
	}
}
func printGameFrame() {
	gameFrameTopLeftRow, gameFrameTopLeftCol := getGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, heigh := GameFrameWidth+2, GameFrameHight+2
	printFilledRect(row, col, width, heigh, GameFrameSymbol)
}

func printFilledRect(row, col, with, heigh int, ch rune) {
	for c := 0; c < with; c++ {
		screen.SetContent(col+c, row, ch, nil, tcell.StyleDefault)
	}

	for r := 1; r < heigh-1; r++ {
		for c := 0; c < with; c++ {
			screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
			screen.SetContent(col+with-1, row+r, ch, nil, tcell.StyleDefault)
		}
	}

	for c := 0; c < with; c++ {
		screen.SetContent(col+c, row+heigh-1, ch, nil, tcell.StyleDefault)
	}

}

func printSnake() {
	for _, p := range snake.parts {
		printFilledRectInGameFrame(p.row, p.col, 1, 1, SnakeSymbol)
		pointsToClear = append(pointsToClear, p)
	}
}

func printFilledRectInGameFrame(row, col, width, heigh int, ch rune) {
	r, c := screenHeigh/2-GameFrameHight/2, screenWidth/2-GameFrameWidth/2
	printFilledRect(row+r, col+c, width, heigh, ch)
}

func printApple() {
	printFilledRectInGameFrame(apple.point.row, apple.point.col, 1, 1, AppleSymbol)
	pointsToClear = append(pointsToClear, apple.point)
}

func getGameFrameTopLeft() (int, int) {
	return screenHeigh/2 - GameFrameHight/2, screenWidth/2 - GameFrameWidth/2
}

func waitForEnter() bool {
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				return true
			} else if ev.Key() == tcell.KeyEscape {
				return false
			}
		}
	}
}
