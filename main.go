package main

import (
	"snake/game"
)

const (
	defaultLength = 10
	defaultHeight = 10
)

func main() {
	gameSnake := game.InitGame(defaultLength, defaultHeight)
	gameSnake.StartGame()
}
