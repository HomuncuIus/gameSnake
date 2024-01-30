package main

import (
	"snake/game"
)

const (
	defaultLength = 15
	defaultHeight = 15
)

func main() {
	gameSnake := game.InitGame(defaultLength, defaultHeight)
	gameSnake.StartGame()
}
