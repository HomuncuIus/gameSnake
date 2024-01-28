package game

type pos struct {
	x, y int
}

func getValue(p *pos) (int, int) {
	return p.x, p.y
}
