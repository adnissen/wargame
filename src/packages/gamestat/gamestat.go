package gamestat

import (
	"reflect"

	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/gameclient"
	"github.com/adnissen/wargame/src/packages/gamemap"
	"github.com/satori/go.uuid"
)

type GameStat struct {
	Armies      []army.Army
	Players     []*gameclient.GameClient
	Uid         uuid.UUID
	Status      string
	CurrentTurn *gameclient.GameClient
	Map         gamemap.Map
}

func (g *GameStat) SendMessageToAllPlayers(mt string, m []byte) {
	for k, _ := range g.Players {
		g.Players[k].SendMessageOfType(mt, m)
	}
}

func (g *GameStat) SetCurrentGameForAllPlayers() {
	for k, _ := range g.Players {
		g.Players[k].SetCurrentGame(g.Uid)
	}
}

func (g *GameStat) EndGame() {
	for k, _ := range g.Players {
		if g.Players[k].IsStillConnected() {
			g.Players[k].ResetCurrentGame()
		}
	}
	g.Status = "ENDED"
}

func CreateGame(p1 *gameclient.GameClient, p2 *gameclient.GameClient) *GameStat {
	if !reflect.DeepEqual(p1.CurrentGame, uuid.UUID{}) || !reflect.DeepEqual(p2.CurrentGame, uuid.UUID{}) {
		return nil
	}
	pary := []*gameclient.GameClient{p1, p2}
	aary := []army.Army{p1.Army, p2.Army}
	gstat := GameStat{Armies: aary, Players: pary, Uid: uuid.NewV4(), Map: gamemap.GetMap()}
	gstat.SetCurrentGameForAllPlayers()
	gstat.SendMessageToAllPlayers("announce", []byte("Game "+gstat.Uid.String()+" starting!"))

	return &gstat
}
