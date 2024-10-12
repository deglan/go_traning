package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF

const PaddleHeigh = 4
const InitialBallVelocityRow = 1
const InitialBallVelocityCol = 2

type GameObject struct {
	row, col, width, heigh int
	velRow, velCol         int
	symbol                 rune
}

var screenWidth, screenHeigh int
var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject

var isGamePaused bool
var debugLog string

var gameObjects = []*GameObject{}

func main() {
	initScreen()
	screenWidth, screenHeigh = screen.Size()

	printIntroductionMessage()
	waitForEnter()
	initGameState()
	inputChan := initUserInput()

	for !isGameOver() {
		handleUserInput(readInput(inputChan))
		updateState()
		drawState()
		time.Sleep(75 * time.Millisecond)
	}

	printWinnerMessage()
}

func updateState() {
	if isGamePaused {
		printString(screenHeigh/2-1, screenWidth/2, "Paused")
		return
	}
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if collidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if collidesWithPaddle(ball, player1Paddle) || collidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}
}

func drawState() {
	if isGamePaused {
		printString(screenHeigh/2-1, screenWidth/2, "Paused")
		return
	}
	screen.Clear()
	printString(0, 0, debugLog)
	for _, obj := range gameObjects {
		print(obj.row, obj.col, obj.width, obj.heigh, obj.symbol)
	}
	screen.Show()
}

func collidesWithWall(obj *GameObject) bool {
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeigh
}

func collidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}
	return collidesOnColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.heigh
}

func isGameOver() bool {
	return getWinner() != ""
}

func getWinner() string {
	if ball.col < 0 {
		return "Player 1"
	} else if ball.col >= screenWidth {
		return "Player 2"
	}
	return ""
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

func initGameState() {
	paddleStart := screenHeigh/2 - PaddleHeigh/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, heigh: PaddleHeigh,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	player2Paddle = &GameObject{
		row: paddleStart, col: screenWidth - 1, width: 1, heigh: PaddleHeigh,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	ball = &GameObject{
		row: screenHeigh / 2, col: screenWidth / 2, width: 1, heigh: 1,
		velRow: InitialBallVelocityRow, velCol: InitialBallVelocityCol,
		symbol: BallSymbol,
	}
	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}
}

func handleUserInput(key string) {
	if key == "Up" && isPlayerInBoundaries(player2Paddle, -1) {

		player2Paddle.row--
	} else if key == "Down" && isPlayerInBoundaries(player2Paddle, 1) {

		player2Paddle.row++
	} else if key == "Rune[w]" && isPlayerInBoundaries(player1Paddle, -1) {

		player1Paddle.row--
	} else if key == "Rune[s]" && isPlayerInBoundaries(player1Paddle, 1) {

		player1Paddle.row++
	} else if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}

}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventResize:
				screen.Sync()
				drawState()

			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}
func isPlayerInBoundaries(player *GameObject, direction int) bool {
	if direction < 0 && player.row <= 0 {
		return false
	}
	if direction > 0 && player.row+player.heigh >= screenHeigh {
		return false
	}
	return true
}

func readInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
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

func print(row, col, width, heigh int, ch rune) {
	for r := 0; r < heigh; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func printIntroductionMessage() {
	screen.Clear()
	printString(screenHeigh/2-1, screenWidth/2, "Hello in Pong Game")
	printString(screenHeigh/2, screenWidth/2, "Click Enter to start the game")
	screen.Show()
}

func printWinnerMessage() {
	winner := getWinner()
	printString(screenHeigh/2-1, screenWidth/2, winner+" won the game!")
	printString(screenHeigh/2, screenWidth/2, "Click Enter to exit")
	screen.Show()
	waitForEnter()
	screen.Fini()
}
