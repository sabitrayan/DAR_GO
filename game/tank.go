package game

import (
	tl "github.com/JoelOtter/termloop"
	"strings"
)

type Tank struct {
	*tl.Entity
	direction int
	color tl.Attr
}

var (
	h = 5
	w = 10
)

const (
	UP    int = 1
	DOWN  int = 2
	LEFT  int = 3
	RIGHT int = 4
)


func (tank *Tank) Draw(screen *tl.Screen) {

	tank.Entity.Draw(screen)
}

func (tank *Tank) Tick(event tl.Event) {}


func TankUpCanvas(tankBodyCell tl.Cell) tl.Canvas {
	canv :=
	"    ██    \n"+
	"██  ██  ██\n"+
	"██████████\n"+
	"██████████\n"+
	"██      ██\n"
	lines := strings.Split(canv, "\n")
	runes := make([][]rune, len(lines))
	for i := range lines {
		runes[i] = []rune(lines[i])
	}
	canvasUp := tl.NewCanvas(w,h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if runes[j][i] == '█'{
				canvasUp[i][j] = tankBodyCell
			}
		}
	}

	return canvasUp
}

func TankDownCanvas(tankBodyCell tl.Cell) tl.Canvas {
	canv := "██      ██\n"+
			"██████████\n"+
			"██████████\n"+
			"██  ██  ██\n"+
			"    ██    "
	lines := strings.Split(canv, "\n")
	runes := make([][]rune, len(lines))
	for i := range lines {
		runes[i] = []rune(lines[i])
	}
	canvasDown := tl.NewCanvas(w,h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if runes[j][i] == '█'{
				canvasDown[i][j] = tankBodyCell
			}
		}
	}

	return canvasDown
}

func TankLeftCanvas(tankBodyCell tl.Cell) tl.Canvas {

	canv :=
		"  ████████\n"+
		"    ████  \n"+
		"████████  \n"+
		"    ████  \n"+
		"  ████████"

	lines := strings.Split(canv, "\n")
	runes := make([][]rune, len(lines))
	for i := range lines {
		runes[i] = []rune(lines[i])
	}
	canvasLeft := tl.NewCanvas(w,h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if runes[j][i] == '█'{
				canvasLeft[i][j] = tankBodyCell
			}
		}
	}

	return canvasLeft
}

func TankRightCanvas(tankBodyCell tl.Cell) tl.Canvas {

	canv :=
		"████████  \n"+
		"  ████    \n"+
		"  ████████\n"+
		"  ████    \n"+
		"████████  "

	lines := strings.Split(canv, "\n")
	runes := make([][]rune, len(lines))
	for i := range lines {
		runes[i] = []rune(lines[i])
	}
	canvasRight := tl.NewCanvas(w,h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if runes[j][i] == '█'{
				canvasRight[i][j] = tankBodyCell
			}
		}
	}



	return canvasRight
}

func init() {

	//	tankBodyCell = tl.Cell{Fg: tl.ColorRed, Color: tl.ColorRed}

	// new blank canvas

}

func NewTankXY(x, y int, cell tl.Cell, direction int) *Tank {

	tank := Tank{
		Entity: tl.NewEntity(x, y, w, h),
	}
	tank.color = cell.Bg
	switch direction {
	case UP:
		TankUp(&tank)
	case RIGHT:
		TankRight(&tank)
	case LEFT:
		TankLeft(&tank)
	case DOWN:
		TankDown(&tank)
	}

	return &tank

}

func TankUp(tank *Tank) {
	canvas := TankUpCanvas(tl.Cell{Bg:tank.color})
	tank.SetCanvas(&canvas)
	tank.direction = UP
}

func TankDown(tank *Tank) {

	canvas := TankDownCanvas(tl.Cell{Bg:tank.color})
	tank.SetCanvas(&canvas)
	tank.direction = DOWN
}

func TankLeft(tank *Tank) {

	canvas := TankLeftCanvas(tl.Cell{Bg:tank.color})
	tank.SetCanvas(&canvas)
	tank.direction = LEFT
}

func TankRight(tank *Tank) {

	canvas := TankRightCanvas(tl.Cell{Bg:tank.color})
	tank.SetCanvas(&canvas)
	tank.direction = RIGHT
}

func (tank *Tank) GetDirection() int {
	return tank.direction
}
