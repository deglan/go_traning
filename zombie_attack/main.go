package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

const GameFrameWidth = 80
const GameFrameHight = 20
const GameFrameSymbol = 0x2588

type Point struct {
	row, col int
	symbol   rune
}

type GameObject struct {
	points         []*Point
	velRow, velCol int
}

var screen tcell.Screen
var screenWidth, screenHeigh int
var isGamePaused bool
var isGameOver bool
var debugLog string
var score int

var player *GameObject
var zombies []*GameObject
var bullets []*GameObject

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
	gameOver()
	screen.Show()
	waitForEnter()
	screen.Fini()
}

func updateState() {
	if isGamePaused {
		return
	}

	moveGameObjects(append(append(zombies, bullets...), player))
	updateZombies()
	collisionDetection()
}

func updateZombies() {
	spawnChange := rand.Intn(100)
	if spawnChange < 5 {
		spawnZombie()
	}
}

func moveGameObjects(obj []*GameObject) {
	for _, obj := range obj {
		for i := range obj.points {
			obj.points[i].row += obj.velRow
			obj.points[i].col += obj.velCol
		}
	}
}

func collisionDetection() {
	objectOutOfBoundsCollision(zombies, true, func(idx int) {
		isGameOver = true
	})

	objectOutOfBoundsCollision(bullets, false, func(idx int) {
		bullets = append(bullets[:idx], bullets[idx+1:]...)
	})

	for _, z := range zombies {
		if areObjectsCollided(player, z, 1) {
			isGameOver = true
		}
	}

	for bi, b := range bullets {
		for zi, z := range zombies {
			if areObjectsCollided(b, z, 1) {
				score++
				bullets = append(bullets[:bi], bullets[bi+1:]...)
				zombies = append(zombies[:zi], zombies[zi+1:]...)
				break
			}
		}
	}
}

// zaczac od 12
func objectOutOfBoundsCollision(objs []*GameObject, lookAhead bool, callback func(int)) {
	for i, obj := range objs {
		velRow, velCol := obj.velRow, obj.velCol
		if lookAhead {
			velRow = 0
			velCol = 0
		}
		if isOutOfBounds(obj, velRow, velCol) {
			callback(i)
			return
		}
	}
}

func areObjectsCollided(obj1, obj2 *GameObject, radius int) bool {
	for _, p1 := range obj1.points {
		for _, p2 := range obj2.points {
			if p1.row == p2.row &&
				math.Abs(float64(p1.col-p2.col)) <= float64(radius) {
				return true
			}
		}
	}
	return false
}

func drawState() {
	debugLog = fmt.Sprintf("Score: %d", score)
	if isGamePaused {
		return
	}
	screen.Clear()
	printString(0, 0, debugLog)
	printGameFrame()
	printGameObject(append(append(zombies, bullets...), player))

	screen.Show()
}

func printGameObject(obj []*GameObject) {
	for _, obj := range obj {
		for _, p := range obj.points {
			printFilledRectInGameFrame(p.row, p.col, 1, 1, p.symbol)
		}
	}
}

func initGameState() {

	player = &GameObject{
		points: []*Point{
			{row: 5, col: 1, symbol: '0'},
			{row: 6, col: 1, symbol: '|'},
			{row: 6, col: 2, symbol: '-'},
			{row: 6, col: 3, symbol: '-'},
			{row: 6, col: 4, symbol: '-'},
			{row: 7, col: 2, symbol: '/'},
			{row: 7, col: 1, symbol: '|'},
			{row: 8, col: 0, symbol: '/'},
			{row: 8, col: 2, symbol: '\\'},
		},
	}
}

func handleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if key == "Enter" {
		spawnBullet(player.points[0].row+1, player.points[0].col+3)
	} else if key == "Rune[w]" && !isOutOfBounds(player, -1, 0) {
		movePlayer(-1, 0)
	} else if key == "Rune[s]" && !isOutOfBounds(player, 1, 0) {
		movePlayer(1, 0)
	} else if key == "Rune[a]" && !isOutOfBounds(player, 0, -1) {
		movePlayer(0, -1)
	} else if key == "Rune[d]" && !isOutOfBounds(player, 0, 1) {
		movePlayer(0, 1)
	}
}

func movePlayer(vRow, vCol int) {
	// this for can modify slice of points
	for i := range player.points {
		player.points[i].row += vRow
		player.points[i].col += vCol
	}
}

func spawnBullet(row, col int) {
	bullets = append(bullets, &GameObject{
		points: []*Point{
			{row: row, col: col, symbol: '*'},
		},
		velRow: 0, velCol: 2,
	})
}

func spawnZombie() {
	originRow, originCol := rand.Intn(GameFrameHight-3), GameFrameWidth-2
	zombies = append(zombies, &GameObject{
		points: []*Point{
			{row: originRow, col: originCol, symbol: '0'},
			{row: originRow + 1, col: originCol, symbol: '|'},
			{row: originRow + 1, col: originCol - 1, symbol: '\\'},
			{row: originRow + 2, col: originCol, symbol: '|'},
			{row: originRow + 3, col: originCol - 1, symbol: '/'},
			{row: originRow + 3, col: originCol + 1, symbol: '\\'},
		},
		velRow: 0, velCol: -1,
	})
}

func isOutOfBounds(obj *GameObject, velRow, velCol int) bool {
	//this for can't modify slice of points
	for _, p := range obj.points {
		targetRow, targetCol := p.row+velRow, p.col+velCol
		if targetRow < 0 ||
			targetRow >= GameFrameHight ||
			targetCol < 0 ||
			targetCol >= GameFrameWidth {
			return true
		}
	}
	return false
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

func printStringCentred(row, col int, str string) {
	col = col - len(str)/2
	printString(row, col, str)
}

func printFilledRectInGameFrame(row, col, width, heigh int, ch rune) {
	r, c := screenHeigh/2-GameFrameHight/2, screenWidth/2-GameFrameWidth/2
	printFilledRect(row+r, col+c, width, heigh, ch)
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

func getGameFrameTopLeft() (int, int) {
	return screenHeigh/2 - GameFrameHight/2, screenWidth/2 - GameFrameWidth/2
}

func gameOver() {
	printStringCentred(screenHeigh/2, screenWidth/2, "Game Over")
	printStringCentred(screenHeigh/2+1, screenWidth/2, fmt.Sprintf("Your score is %d", score))
	printStringCentred(screenHeigh/2+2, screenWidth/2, "Pres enter to quit")
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
