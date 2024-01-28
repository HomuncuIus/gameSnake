package game

import (
	"container/list"
	"math/rand"
	"reflect"
)

const (
	snakeIcon          = "o"
	snakeHeadUpIcon    = "^"
	snakeHeadDownIcon  = "v"
	snakeHeadLeftIcon  = "<"
	snakeHeadRightIcon = ">"
)
const (
	up = iota
	down
	left
	right
)

type snake struct {
	locations *list.List
	direction int
}

func (s *snake) getHead() *pos {
	p, ok := s.locations.Front().Value.(pos)
	if ok {
		return &p
	}
	return nil
}

func (s *snake) getTail() *pos {
	p, ok := s.locations.Back().Value.(pos)
	if ok {
		return &p
	}
	return nil
}

func (s *snake) updateLocation() (bool, interface{}) {
	headPos := s.getHead()
	var frontPos pos
	switch s.direction {
	case up:
		frontPos = pos{headPos.x, headPos.y - 1}
	case down:
		frontPos = pos{headPos.x, headPos.y + 1}
	case left:
		frontPos = pos{headPos.x - 1, headPos.y}
	case right:
		frontPos = pos{headPos.x + 1, headPos.y}
	default:
		return false, nil
	}
	if !s.validateBody(&frontPos) {
		return false, nil
	}
	s.locations.PushFront(frontPos)
	return true, s.locations.Remove(s.locations.Back())
}

func (s *snake) growUp(p pos) {
	s.locations.PushBack(p)
}

func (s *snake) validateBody(p *pos) bool {
	for loc := s.locations.Front(); loc != nil; loc = loc.Next() {
		bodyPos := loc.Value.(pos)
		if reflect.DeepEqual(p, &bodyPos) {
			return false
		}
	}
	return true
}

func (s *snake) changeDirection(direction int) {
	s.direction = direction
}

func initSnake(length, height int) *snake {
	x, y := rand.Intn(length/2)+length/4, rand.Intn(height/2)+height/4
	location := list.New()
	location.PushFront(pos{x, y})
	return &snake{location, up}
}
