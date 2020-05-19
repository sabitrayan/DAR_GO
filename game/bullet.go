package game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	speed float32 = 0.3
)

type Bullet struct {
	*tl.Entity
	direction int
	owner string
	x,y float32
}

func NewBullet(x, y, d int, owner string) *Bullet {
	b := Bullet{
		Entity:    tl.NewEntity(x, y, 1, 1),
		direction: d,
		owner: 		owner,
		x:float32(x),
		y:float32(y),
	}
	b.SetCell(0, 0, &tl.Cell{Fg: tl.ColorBlack, Bg: tl.ColorBlack})

	return &b
}

func (b *Bullet) Draw(screen *tl.Screen) {
	bX, bY := b.Position()
	screenX, screenY := Game().Size()
	if bX > screenX || bX < 0 || bY > screenY || bY < 0 {
		screen.RemoveEntity(b)
		screen.Level().RemoveEntity(b)
	}
	switch b.direction {

	case UP:
		b.y-=speed
	case DOWN:
		b.y+=speed
	case LEFT:
		b.x-=2.0*speed
	case RIGHT:
		b.x+=2.0*speed
	}
	b.SetPosition(int(b.x), int(b.y))
	b.Entity.Draw(screen)

}

func (b *Bullet) Tick(event tl.Event) {
	for id, p := range Game().onlinePlayers {
		px, py := p.Position()
		pw, ph := p.Size()

		bx, by := b.Position()
		bw, bh := b.Size()

		if px < bx+bw && px+pw > bx &&
			py < by+bh && py+ph > by {
			if id != Game().currentPlayerID && b.owner == Game().currentPlayerID{
				if p.HP > 5 {
					Game().onlinePlayers[id].HP -= 5
					Command{ID: id, ReplyTo: b.owner,Action: ATTACKED}.Send()
				}else{
					Command{ID: p.ID, Action: DELETE}.Send()
					Command{ID: b.owner, ReplyTo: b.owner, Action:KILL}.Send()

					Game().onlinePlayers[b.owner].Score += 1
					Game().onlinePlayers[b.owner].HP += 50

					playersToDelete[p.ID] = true
					Game().level.RemoveEntity(Game().onlinePlayers[p.ID])
					delete(Game().onlinePlayers, p.ID)
				}
			}
			Game().Level().RemoveEntity(b)
		}
	}
}

func (b *Bullet) Collide(p tl.Physical) {
	if _, ok := p.(*tl.Rectangle); ok {
		if p.(*tl.Rectangle).Color() != avoidanceBorderColor{
			Game().Level().RemoveEntity(b)
		}
	}
}

func (b *Bullet) Position() (int,int){
	return int(b.x), int(b.y)
}

