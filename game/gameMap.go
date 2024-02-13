package game

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	gameMapIcon           = "."
	initialDurationFactor = 1000
	maxDurationFactor     = 200
)

var seed = time.Now().UnixNano()

type Game struct {
	length, height int
	panel          [][]string
	snake          *snake
	bean           *bean
	score          int
	gameOver       bool
}

// display plant the game map with snake and bean on
func (g *Game) display() {
	for i := 0; i <= g.length; i++ {
		fmt.Print("---")
	}
	fmt.Println("-")
	for y := 0; y < g.height; y++ {
		fmt.Print("|  ")
		for x := 0; x < g.length; x++ {
			fmt.Print(g.panel[y][x] + "  ")
		}
		fmt.Println("|")
	}
	for i := 0; i <= g.length; i++ {
		fmt.Print("---")
	}
	fmt.Println("-")
}

// updateItems check whether the snake eats the bean and generate a new bean if it does
func (g *Game) updateItems(snakeLastPos *pos) {
	snakeHead := g.snake.getHead()
	if !reflect.DeepEqual(snakeHead, g.bean.location) {
		return
	}
	g.snake.growUp(*snakeLastPos)
	g.score += 1
	g.generateBean()
}

// generateBean generate a new bean at the position where snake not occupies
func (g *Game) generateBean() {
	occupiedPos := make(map[pos]interface{})
	for h := 0; h < g.height; h++ {
		for l := 0; l < g.length; l++ {
			occupiedPos[pos{l, h}] = nil
		}
	}
	for loc := g.snake.units.Front(); loc != nil; loc = loc.Next() {
		pos := loc.Value.(pos)
		delete(occupiedPos, pos)
	}
	restPosNum := len(occupiedPos)
	choice := rand.Intn(restPosNum)
	for pos := range occupiedPos {
		if choice != 0 {
			choice -= 1
		}
		g.bean = &bean{&pos}
		return
	}
}

// updatePanel update the game panel
func (g *Game) updatePanel() {
	beanX, beanY := getValue(g.bean.location)
	g.panel[beanY][beanX] = beanIcon
	for loc := g.snake.units.Front(); loc != nil; loc = loc.Next() {
		pos := loc.Value.(pos)
		g.panel[pos.y][pos.x] = snakeIcon
	}
	snakeHeadX, snakeHeadY := getValue(g.snake.getHead())
	switch g.snake.direction {
	case up:
		g.panel[snakeHeadY][snakeHeadX] = snakeHeadUpIcon
	case down:
		g.panel[snakeHeadY][snakeHeadX] = snakeHeadDownIcon
	case left:
		g.panel[snakeHeadY][snakeHeadX] = snakeHeadLeftIcon
	case right:
		g.panel[snakeHeadY][snakeHeadX] = snakeHeadRightIcon
	}

}

// clearLastPos clean the last footprint of snake
func (g *Game) clearLastPos(snakeLastPos *pos) {
	lastX, lastY := getValue(snakeLastPos)
	g.panel[lastY][lastX] = gameMapIcon

}

// checkBoundary check whether the snake hit the map boundary
func (g *Game) checkBoundary() bool {
	headPos := g.snake.getHead()
	if headPos == nil {
		return false
	}
	if headPos.x < 0 || headPos.x >= g.length {
		return false
	}
	if headPos.y < 0 || headPos.y >= g.height {
		return false
	}
	return true
}

func (g *Game) getDurationFactor() int {
	return maxInt(maxDurationFactor, initialDurationFactor-g.score*10)
}

func (g *Game) StartGame() {
	ticker := time.NewTicker(initialDurationFactor * time.Millisecond)
	var err error
	err = keyboard.Open()
	defer func() {
		err = keyboard.Close()
		if err != nil {
			fmt.Printf("an error occerd: %v \n", err)
		}
	}()

	go keyInput(g)

	for {
		select {
		case <-ticker.C:
			ticker.Reset(time.Duration(g.getDurationFactor()) * time.Millisecond)
			err = clearOutput()
			if err != nil {
				g.gameOver = true
			}
			updateSuccess, lastPos := g.snake.updateLocation()
			if !updateSuccess {
				g.gameOver = true
			}
			if !g.checkBoundary() {
				g.gameOver = true
			}
			if g.gameOver {
				fmt.Printf("Game is over, Your score is %v\n", g.score)
				return
			}
			tailPos := lastPos.(pos)
			g.clearLastPos(&tailPos)
			g.updateItems(&tailPos)
			g.updatePanel()
			g.display()
		}
	}
}

func InitGame(length, height int) *Game {
	rand.Seed(seed)
	panel := make([][]string, height)
	for y := range panel {
		panel[y] = make([]string, length)
		for x := range panel[y] {
			panel[y][x] = gameMapIcon
		}
	}
	snake := initSnake(length, height)
	game := &Game{length, height, panel, snake, nil, 0, false}
	game.generateBean()
	return game
}

func clearOutput() error {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	err := c.Run()
	return err
}

func keyInput(g *Game) {
	go func() {
		for {
			_, key, _ := keyboard.GetKey()
			switch key {
			case keyboard.KeyArrowUp:
				g.snake.changeDirection(0)
			case keyboard.KeyArrowDown:
				g.snake.changeDirection(1)
			case keyboard.KeyArrowLeft:
				g.snake.changeDirection(2)
			case keyboard.KeyArrowRight:
				g.snake.changeDirection(3)
			default:
			}
		}
	}()
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
