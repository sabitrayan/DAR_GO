package game

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
	"math"
	"sort"
	"strconv"
)

var (
	menuHidden      = false
	nameOnChange    = false
	currentUsername string
	newUsername     []rune
)

type menu struct {
	items []*tl.Text
	index int
}

const (
	START_OR_RESUME_GAME = 0
	CHANGENAME 	= 	1
	EXIT 		=	2
)

func Username() string {
	return currentUsername
}

func SetUsername(username string) {
	currentUsername = username
}

func (m *menu) Tick(ev tl.Event) {
	if _, prs := Game().onlinePlayers[Game().currentPlayerID]; !prs {
		m.items[0].SetText("Start game")
		menuHidden = false
	}else{
		m.items[0].SetText("Resume")
	}
	if ev.Type == tl.EventKey {
		if !nameOnChange {
			switch ev.Key{
			case tl.Key(65535):
				menuHidden = false
				m.index = 0
			case tl.KeyArrowUp:
				if m.index > 0 && !menuHidden {
					m.index -= 1
				}
			case tl.KeyArrowDown:
				if m.index < len(m.items)-1 && !menuHidden {
					m.index += 1
				}
			case tl.KeyEnter:
				switch m.index {
				case START_OR_RESUME_GAME:
					if _, prs := Game().onlinePlayers[Game().currentPlayerID]; !prs {
						NewPlayer()
					}
					menuHidden = true
				case CHANGENAME:
					nameOnChange = true
					newUsername = make([]rune, 0)
				case EXIT:
					Command{ID: Game().currentPlayerID, Action:DELETE}.Send()
					Game().Stop()
				}
			}
		} else {
			switch  ev.Key {
			case tl.Key(65535) :
				nameOnChange = false
			case tl.Key(127):
				if len(newUsername) > 0 {
					newUsername = newUsername[:len(newUsername)-1]
				}
			case tl.KeyEnter:
				if len(newUsername) > 0{
					SetUsername(string(newUsername))
				}
				nameOnChange = false
			default:
				if ((ev.Ch >= '0' && ev.Ch <= '9' ) || (ev.Ch >= 'a' && ev.Ch <= 'z' )|| (ev.Ch >= 'A' && ev.Ch <= 'Z') ||
					(ev.Ch >= 'а' && ev.Ch <= 'я' )|| (ev.Ch >= 'А' && ev.Ch <= 'Я') ) && len(newUsername) < 15{
					newUsername = append(newUsername, ev.Ch)
				}
			}
		}
	}
}
func (m *menu) Draw(s *tl.Screen) {


	if !menuHidden {
		Game().Level().SetOffset(0,0)
		sx, sy := s.Size()
		if nameOnChange {

			tl.NewText(sx/2-11, 0, "press F1 to GO the BACK", tl.ColorRed, tl.ColorYellow).Draw(s)

			text := tl.NewText(0,0, "Enter username(at most 15 chars):", tl.ColorBlack, tl.ColorWhite)
			ix, iy := text.Size()
			text.SetPosition(sx/2-ix/2, sy/2-iy/2)
			text.Draw(s)

			username := tl.NewText(0,0, string(newUsername), tl.ColorBlack, tl.ColorWhite)
			ix, iy = username.Size()
			username.SetPosition(sx/2-ix/2, sy/2-iy/2+1)
			username.Draw(s)

		} else {
			for i,_ := range m.items {
				if i == m.index {
					m.items[i].SetColor(tl.ColorWhite, tl.ColorBlack)
				} else {
					m.items[i].SetColor(tl.ColorBlack, tl.ColorWhite)
				}

				ix, iy := m.items[i].Size()
				m.items[i].SetPosition(sx/2-ix/2, sy/2-iy/2+i-len(m.items))
				m.items[i].Draw(s)
			}
		}
	} else {
		top := []struct{
			id string
			score int
			hp int
			username string
		}{}
		for _, p := range Game().onlinePlayers {
			top = append(top, struct {
				id string
				score    int
				hp       int
				username string
			}{id: p.ID,score: p.Score, hp: p.HP, username: p.Username})
		}

		// sorting comparator
		sort.Slice(top, func(i, j int) bool {
			if top[i].score > top[j].score {
				return true
			}else if top[i].score < top[j].score {
				return false
			}else if top[i].hp > top[j].hp {
				return true
			}else if top[i].hp < top[j].hp {
				return false
			}else if top[i].username < top[j].username {
				return true
			}else if top[i].username > top[j].username {
				return false
			}else{
				return top[i].id < top[j].id
			}
		})

		x,y := Game().Level().Offset()
		x,y = -x,-y

		tl.NewText(x+1, y, "       TOP 5       ", tl.ColorRed, tl.ColorWhite).Draw(s)
		tl.NewText(x+1, y+1, "# Score HP  Username", tl.ColorGreen, tl.ColorWhite).Draw(s)

		sW, _ := Game().Screen().Size()
		instruction := "--press F1 to GO the MENU--Space to fire--Arrows to move--"
		tl.NewText(x+sW/2-len(instruction)/2, y, instruction, tl.ColorRed, tl.ColorYellow).Draw(s)

		for i:=0; i < int(math.Min(float64(len(top)), 5.0)) ; i++{
			tl.NewText(x+1, y+2+i, fmt.Sprintf("%-2d%-5d %-3d %s",i+1, top[i].score, top[i].hp, top[i].username), Game().onlinePlayers[top[i].id].color, tl.ColorWhite).Draw(s)
		}
	}
}

func StartMenu(){

	arg_items := []string{"Start Game", "Set Name", "Exit"}

	SetUsername("Player_NO_NAME_"+strconv.Itoa(len(Game().onlinePlayers)))

	menu := menu{
		items:    make([]*tl.Text, len(arg_items)),
		index: 	  0,
	}
	for i, item := range arg_items {
		menu.items[i] = tl.NewText(0,i, item, tl.ColorBlack, tl.ColorWhite)
		if i == menu.index {
			menu.items[i] = tl.NewText(0,i, item, tl.ColorWhite, tl.ColorBlack)
		}
	}

	Game().Level().AddEntity(&menu)
}

