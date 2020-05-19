package game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	borderrColor         = tl.Attr(1)
	edgeColor            = tl.ColorBlack
	avoidanceBorderColor = tl.Attr(9)
)


func NewWalls() []*tl.Rectangle {
	walls := make([]*tl.Rectangle, 0)

	gX, gY := Game().Size()



	walls = append(
		walls,
		tl.NewRectangle(10,20,5,25, borderrColor),
		tl.NewRectangle(10,20,20,2, borderrColor),
		tl.NewRectangle(30,22,5,2, borderrColor),
		tl.NewRectangle(30,43,5,2, borderrColor),
		tl.NewRectangle(35,24,5,19, avoidanceBorderColor),
		tl.NewRectangle(10,45,20,2, borderrColor),

		tl.NewRectangle(51,20,5,27, borderrColor),
		tl.NewRectangle(71,20,5,27, borderrColor),
		tl.NewRectangle(56,20,15,2, avoidanceBorderColor),
		tl.NewRectangle(56,30,15,2, avoidanceBorderColor),


		tl.NewRectangle(87,20,5,27, borderrColor),
		tl.NewRectangle(87,20,25,2, borderrColor),
		tl.NewRectangle(107,22,5,10, avoidanceBorderColor),
		tl.NewRectangle(87,32,25,2, borderrColor),

		tl.NewRectangle(92,34,5,2, borderrColor),
		tl.NewRectangle(95,36,5,2, borderrColor),
		tl.NewRectangle(99,38,5,3, borderrColor),
		tl.NewRectangle(103,41,5,3, borderrColor),
		tl.NewRectangle(107,44,5,3, borderrColor),





		tl.NewRectangle(-1,-1,gX+2,1, edgeColor),
		tl.NewRectangle(-1,-1,1,gY+2, edgeColor),
		tl.NewRectangle(gX,-1,1,gY+2, edgeColor),
		tl.NewRectangle(-1,gY,gX+2,1, edgeColor),
	)

	return walls
}