package gamestat

import (
	"encoding/json"
	"fmt"
	"reflect"

	"strconv"

	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/dice"
	"github.com/adnissen/wargame/src/packages/gameclient"
	"github.com/adnissen/wargame/src/packages/gamemap"
	"github.com/adnissen/wargame/src/packages/units"
	"github.com/satori/go.uuid"
)

type GameStat struct {
	Armies           []army.Army
	Players          []*gameclient.GameClient
	Uid              uuid.UUID
	Status           string
	CurrentTurn      int
	Map              gamemap.Map
	UnitActionCounts map[*units.Unit]int
}

func (g *GameStat) SendMessageToAllPlayers(mt string, m []byte) {
	for k, _ := range g.Players {
		g.Players[k].SendMessageOfType(mt, m)
	}
}

func (g *GameStat) ResetActions() {
	if g.UnitActionCounts == nil {
		g.UnitActionCounts = make(map[*units.Unit]int)
	}
	for k, _ := range g.Armies {
		for i, _ := range g.Armies[k].Squads {
			for j, _ := range g.Armies[k].Squads[i].Grunts {
				g.UnitActionCounts[&g.Armies[k].Squads[i].Grunts[j]] = 2
			}
			g.UnitActionCounts[&g.Armies[k].Squads[i].Leader] = 2
		}
	}
}

func (g *GameStat) GetUnitOnTile(x int, y int) *units.Unit {
	return g.Map.Map[x][y].Unit
}

func (g *GameStat) SetCurrentGameForAllPlayers() {
	for k, _ := range g.Players {
		g.Players[k].SetCurrentGame(g.Uid)
	}
}

func (g *GameStat) GetMapJson() []byte {
	j, _ := json.Marshal(g.Map)
	return j
}

func (g *GameStat) GetUnitJson() []byte {
	j, _ := json.Marshal(g.Armies)
	return j
}

func (g *GameStat) EndGame() {
	for k, _ := range g.Players {
		if g.Players[k].IsStillConnected() {
			g.Players[k].ResetCurrentGame()
		}
	}
	g.Status = "ENDED"
}

func (g *GameStat) Attack(attacker *units.Unit, defender *units.Unit) {
	if g.UnitActionCounts[attacker] > 0 {
		if gamemap.DistanceBetweenTiles(attacker.X, attacker.Y, defender.X, defender.Y) <= attacker.Attributes.Rng {
			r := dice.Roll(20)
			g.SendMessageToAllPlayers("announce", []byte(attacker.DisplayName+"("+strconv.Itoa(attacker.Attributes.Atk)+") rolls "+strconv.Itoa(r)+" against "+defender.DisplayName+"("+strconv.Itoa(defender.Attributes.Def)+")"))
			if (r + attacker.Attributes.Atk) > defender.Attributes.Def {
				defender.Attributes.Hps -= attacker.Attributes.Dmg
				if defender.Attributes.Hps <= 0 {
					delete(g.UnitActionCounts, defender)
				}
			}
			//attacking is the last thing you can do
			g.UnitActionCounts[attacker] = 0
		}
	}
}

func (g *GameStat) SpawnAllUnits() {
	for team := range g.Armies {
		fmt.Print("Army team: ")
		fmt.Print(team)
		fmt.Print("\n")
		fmt.Println("==============")
		for s := range g.Armies[team].Squads {
			for k := range g.Armies[team].Squads[s].Grunts {
				//needs to be a pointer to this array? makes changes to this array but loop still uses old copy
				grunt := &g.Armies[team].Squads[s].Grunts[k]
				g.Map.SpawnUnitOnFirstAvailable(grunt, team)
				if grunt.Spawned {
					fmt.Print("spawned unit at ")
					fmt.Print(grunt.X)
					fmt.Print(" ")
					fmt.Print(grunt.Y)
					fmt.Print("\n")
					//fmt.Println(grunt == &g.Armies[team].Squads[s].Grunts[k])
					//fmt.Println(grunt == g.GetUnitOnTile(grunt.X, grunt.Y))
					//fmt.Println(grunt.X == g.GetUnitOnTile(grunt.X, grunt.Y).X && grunt.Y == g.GetUnitOnTile(grunt.X, grunt.Y).Y)
				}
			}
			g.Map.SpawnUnitOnFirstAvailable(&g.Armies[team].Squads[s].Leader, team)
		}
	}
}

func CreateGame(p1 *gameclient.GameClient, p2 *gameclient.GameClient) *GameStat {
	if !reflect.DeepEqual(p1.CurrentGame, uuid.UUID{}) || !reflect.DeepEqual(p2.CurrentGame, uuid.UUID{}) {
		return nil
	}
	pary := []*gameclient.GameClient{p1, p2}
	aary := []army.Army{p1.Army, p2.Army}
	gstat := GameStat{Armies: aary, Players: pary, Uid: uuid.NewV4(), Map: gamemap.GetCustomMap()}
	gstat.SetCurrentGameForAllPlayers()
	gstat.SendMessageToAllPlayers("announce", []byte("Game "+gstat.Uid.String()+" starting!"))
	gstat.SendMessageToAllPlayers("map_data", gstat.GetMapJson())
	fmt.Println("BEFOReTHIGN")
	fmt.Println(&gstat.Armies[0].Squads[0].Grunts[0] == &gstat.Armies[0].Squads[1].Grunts[0])
	gstat.SpawnAllUnits()
	fmt.Println("THIGN")
	fmt.Println(&gstat.Armies[0].Squads[0] == &gstat.Armies[0].Squads[1])
	fmt.Println(&gstat.Armies[0].Squads[0].Grunts[0] == &gstat.Armies[0].Squads[1].Grunts[0])
	gstat.SendMessageToAllPlayers("game_start_army_data", gstat.GetUnitJson())
	gstat.ResetActions()

	return &gstat
}
