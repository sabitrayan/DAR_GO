package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"math/rand"
	"strconv"
	"time"
)
var (
	lastShotTime    = time.Now() // bullet per time
	bulletPerSecond = 3
)



type player struct {
	*Tank
	preX, preY	int
	ID 			string
	Username 	string
	Score		int
	HP 			int

	lastOnlineTime time.Time
}

func Player(id string) *player {
	return &player{
		Tank:     NewTankXY(0,0, tl.Cell{Bg: tl.ColorBlack}, UP),
		ID:       id,
	}
}

func NewPlayer() *player{
	gX, gY := Game().Size()

	R := randInt(0,256)
	G := randInt(0,256-R)
	B := randInt(0,256-G-R)

	player := &player{
		Tank:     NewTankXY(
			randInt(0,gX-9),
			randInt(0, gY-9),
			tl.Cell{Bg:tl.RgbTo256Color(R,G,B)},
			randInt(1,4),
			),
		ID:       genRandString(32),
		Username: Username(),
		Score:    0,
		HP:       100,

	}

	// avoid collision
	for player.CollideWorker(player.direction) {
		player.direction = randInt(1,4)
		player.SetPosition(	randInt(0,gX-9), randInt(0, gY-9))
	}

	player.lastOnlineTime = time.Now()
	Game().onlinePlayers[player.ID] = player
	Game().Screen().Level().AddEntity(player)

	go func() {
		for {
			Publisher(OnlineExchange(), "", amqp.Publishing{Body: []byte(player.ID)})
			time.Sleep(time.Second*2)
		}
	}()

	Game().currentPlayerID = player.ID

	// request for positions of other's tanks
	Command{Action:REQUEST}.Send()

	return player
}

func (p *player) Tick(event tl.Event) {
	if p.ID == Game().currentPlayerID && menuHidden && event.Type == tl.EventKey {
		p.preX, p.preY = p.Position()
		command := Command{ID: p.ID, Action: TANK, X: p.preX, Y: p.preY, Direction: p.direction, Score:-1}

		switch event.Key {
		case tl.KeyArrowUp:
			if p.direction == UP {
				command.Y -= 1
				if p.CollideWorker(UP) {
					return
				}
			}
			command.Direction = UP
		case tl.KeyArrowDown:
			if p.direction == DOWN {
				command.Y += 1
				if p.CollideWorker(DOWN) {
					return
				}
			}
			command.Direction = DOWN
		case tl.KeyArrowRight:
			if p.direction == RIGHT {
				command.X += 2
				if p.CollideWorker(RIGHT) {
					return
				}
			}
			command.Direction = RIGHT
		case tl.KeyArrowLeft:
			if p.direction == LEFT {
				command.X -= 2
				if p.CollideWorker(LEFT) {
					return
				}
			}
			command.Direction = LEFT
		case tl.KeySpace:
			if time.Now().Sub(lastShotTime).Milliseconds() >= int64(1000/bulletPerSecond) {
				var bulletX, bulletY, bulletDirection  int
				bulletDirection = p.Tank.GetDirection()

				switch bulletDirection {
				case UP:
					bulletX = p.preX + 4
					bulletY = p.preY - 1
				case DOWN:
					bulletX = p.preX + 4
					bulletY = p.preY + 6
				case LEFT:
					bulletX = p.preX - 1
					bulletY = p.preY + 2
				case RIGHT:
					bulletX = p.preX + 10
					bulletY = p.preY + 2
				}
				command.Action = BULLET
				command.X = bulletX
				command.Y = bulletY
				command.Direction = bulletDirection
				lastShotTime = time.Now()
			}
		}
		command.Send()
	}
}

func (p *player) Draw(screen *tl.Screen) {
	if menuHidden {
		x,y := p.Position()
		teX := len([]rune(p.Username))
		tl.NewText(x+3,y-2,strconv.Itoa(p.HP), tl.ColorBlue, tl.ColorWhite).Draw(screen)
		tl.NewText(x+4-teX/2,y-1, p.Username, tl.ColorBlack, tl.ColorWhite).Draw(screen)
		p.Entity.Draw(screen)
	}
}

func (p *player) CollideWorker(direction int) bool {
	if p.ID == Game().currentPlayerID {
		for _, c := range Game().onlinePlayers {
			if c.ID != p.ID && collided(p, c, direction) {
				return true
			}
		}
		for _, c := range Game().walls {
			if collided(p, c, direction) {
				return true
			}
		}
	}
	return false
}

func collided(p tl.Physical, c tl.Physical, direction int) bool {
	if _, ok := c.(*tl.Rectangle); ok && c.(*tl.Rectangle).Color() == avoidanceBorderColor {
		return false
	}
	px, py := p.Position()
	cx, cy := c.Position()
	pw, ph := p.Size()
	cw, ch := c.Size()

	switch direction {
	case UP:
		py -= 1
	case DOWN:
		py += 1
	case LEFT:
		px -= 2
	case RIGHT:
		px += 2
	}

	if px < cx+cw && px+pw > cx &&
		py < cy+ch && py+ph > cy {
		return true
	}
	return false
}



func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func genRandString(l int) string {
	bytes := make([]byte, l)

	for i:=0; i < l; i++ {
		bytes[i] = byte(randInt(65,90))
	}
	return string(bytes)
}





