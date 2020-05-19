package game

import (
	"encoding/json"
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"time"
)

var (
	game 	*gameClass
	fps 			= 	90.0 // should be float
	gameWidth 		=	150
	gameHeight 		=	60
	backgroundColor = 	tl.ColorWhite
	playersToDelete = 	make(map[string]bool)
)

type gameClass struct {
	*tl.Game

	width, height 	int
	level 			*tl.BaseLevel

	currentPlayerID string

	walls 			[]*tl.Rectangle
	onlinePlayers 	map[string]*player

}

func Game() *gameClass {
	if game == nil {
		game = &gameClass{
			Game:   tl.NewGame(),
			width:  gameWidth,
			height: gameHeight,
			level: tl.NewBaseLevel(tl.Cell{Bg: backgroundColor}),
			onlinePlayers: make(map[string]*player),
		}
		game.walls = NewWalls()
	}
	return game
}

func (g *gameClass) Size() (w int, h int){
	return g.width, g.height
}

func (g *gameClass) Level() *tl.BaseLevel {
	return g.level
}

func (g *gameClass) Start() {

	g.Screen().SetFps(fps)

	for _, w := range g.walls {
		g.level.AddEntity(w)
	}

	g.Screen().SetLevel(g.level)

	g.getOnlinePlayers()

	g.checkOnlinePlayers()

	g.listenCommands()

	StartMenu()

	g.Game.Start()
}


func (g *gameClass) getOnlinePlayers() {

	msgs := Consumer(OnlineQueue())

	go func() {
		for d := range msgs {
			id := string(d.Body)

			if _, ok := playersToDelete[id]; !ok{
				if _, ok := g.onlinePlayers[id]; !ok {
					g.onlinePlayers[id] = Player(id)
					game.Screen().Level().AddEntity(g.onlinePlayers[id])
					Command{Action:REQUEST, ReplyTo:ReceiverQueue()}.Send()
				}
				g.onlinePlayers[id].lastOnlineTime = time.Now()
			}

			d.Ack(false)
		}
	}()
}

func (g *gameClass) checkOnlinePlayers() {

	go func() {
		for {
			for id,p := range g.onlinePlayers {
				if time.Now().Sub(p.lastOnlineTime).Seconds() > 10 {
					Command{ID:id, Action:DELETE,}.Send()
				}
			}
			time.Sleep(time.Second)
		}
	}()
}

func (g *gameClass) listenCommands() 	{

	msgs := Consumer(ReceiverQueue())

	go func() {

		for d := range msgs {
			c := Command{}
			err := json.Unmarshal(d.Body, &c)
			failOnError(err, "Failed to unmarshal command", "")

			switch c.Action {
			case KILL:
				if _, ok := g.onlinePlayers[c.ID]; ok && c.ReplyTo != g.currentPlayerID{
					g.onlinePlayers[c.ID].HP += 50
					g.onlinePlayers[c.ID].Score += 1
				}
			case ATTACKED:
				if _, ok := g.onlinePlayers[c.ID]; ok && c.ReplyTo != g.currentPlayerID{
					g.onlinePlayers[c.ID].HP -= 5
				}
			case TANK:
				if _, ok := g.onlinePlayers[c.ID]; ok {
					p := g.onlinePlayers[c.ID]

					p.Username = c.Username
					if c.Score != -1 {
						p.Score = c.Score
						p.HP = c.HP
						p.color = c.Color
					}
					p.preX, p.preY = p.Position()
					p.SetPosition(c.X, c.Y)
					switch c.Direction {
					case UP:
						TankUp(p.Tank)
					case DOWN:
						TankDown(p.Tank)
					case LEFT:
						TankLeft(p.Tank)
					case RIGHT:
						TankRight(p.Tank)
					}
					if p.ID == Game().currentPlayerID && menuHidden {
						sW, sH := Game().Screen().Size()
						tX, tY := p.Position()
						padding := 2
						sX := min(gameWidth-sW+padding, max(-padding,tX-sW/2+5))
						sY := min(gameHeight-sH+padding,max(-padding,tY-sH/2+5))
						Game().Level().SetOffset(-sX, -sY)
					}

				}
			case BULLET:
				b := NewBullet(c.X, c.Y, c.Direction, c.ID)
				g.level.AddEntity(b)
			case DELETE:
				if _, ok := g.onlinePlayers[c.ID]; ok {

					playersToDelete[c.ID] = true
					g.level.RemoveEntity(g.onlinePlayers[c.ID])
					delete(g.onlinePlayers, c.ID)
				}
			case REQUEST:
				if _, ok := g.onlinePlayers[g.currentPlayerID]; ok {
					player := g.onlinePlayers[g.currentPlayerID]
					x,y := player.Position()
					Command{ID: player.ID, Action: TANK, X:x, Y:y, Direction:player.GetDirection(), Score: player.Score, ReplyTo: d.CorrelationId}.Send()
				}
			}

			d.Ack(false)
		}
	}()

}


const (
	TANK     = 0
	BULLET   = 1
	DELETE   = 2
	REQUEST  = 3
	ATTACKED = 4
	KILL     = 5
) // command action

type Command struct {
	ID string
	Action int
	ReplyTo string
	X, Y, Direction, Score, HP int
	Username string
	Color tl.Attr
}

func (c Command) Send(){

	if c.Action == TANK {
		p := Game().onlinePlayers[Game().currentPlayerID]
		c.HP = p.HP
		//c.Score = p.Score
		c.Username = Username()
		c.Color = p.color
	}

	body, _ := json.Marshal(c)

	msg := amqp.Publishing{
		ContentType:     "application/json",
		Body:            body,
	}
	if c.Action == TANK && c.ReplyTo != "" {
		Publisher("",c.ReplyTo, msg)
	}else{
		Publisher(CommandsExchange(),"", msg)
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}