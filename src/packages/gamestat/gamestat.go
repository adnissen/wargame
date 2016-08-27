package gamestat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"strconv"

	"regexp"

	"github.com/adnissen/go.uuid"
	"github.com/adnissen/gorm"
	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/dice"
	"github.com/adnissen/wargame/src/packages/gameclient"
	"github.com/adnissen/wargame/src/packages/gamemap"
	"github.com/adnissen/wargame/src/packages/los"
	"github.com/adnissen/wargame/src/packages/units"
)

type GameStat struct {
	Armies           []army.Army
	Players          []*gameclient.GameClient
	Points           []int
	Uid              uuid.UUID
	Status           string
	CurrentTurn      int
	Map              gamemap.Map
	UnitActionCounts map[string]int
	UnitCombatCounts map[string]int
}

func (g *GameStat) SendMessageToAllPlayers(mt string, m []byte) {
	for k, _ := range g.Players {
		g.Players[k].SendMessageOfType(mt, m)
	}
}

func (g *GameStat) SendMessageToPlayer(p *gameclient.GameClient, mt string, m []byte) {
	for k, _ := range g.Players {
		if g.Players[k] == p {
			g.Players[k].SendMessageOfType(mt, m)
		}
	}
}

func (g *GameStat) ResetActions() {
	if g.UnitActionCounts == nil {
		g.UnitActionCounts = make(map[string]int)
		g.UnitCombatCounts = make(map[string]int)
	}
	for k, _ := range g.Armies {
		for i, _ := range g.Armies[k].Squads {
			for j, _ := range g.Armies[k].Squads[i].Grunts {
				g.UnitActionCounts[g.Armies[k].Squads[i].Grunts[j].Uid.String()] = 2
				g.UnitCombatCounts[g.Armies[k].Squads[i].Grunts[j].Uid.String()] = 1
			}
			g.UnitActionCounts[g.Armies[k].Squads[i].Leader.Uid.String()] = 2
			g.UnitCombatCounts[g.Armies[k].Squads[i].Leader.Uid.String()] = 1
		}
	}
}

func (g *GameStat) GetPlayerIndex(player *gameclient.GameClient) int {
	for i, p := range g.Players {
		if p == player {
			return i
		}
	}
	return -1
}

func (g *GameStat) GetUnitOnTile(x int, y int) *units.Unit {
	return g.Map.Map[x][y].Unit
}

func (g *GameStat) GetTile(x int, y int) *gamemap.Tile {
	return &g.Map.Map[x][y]
}

func (g *GameStat) SetCurrentGameForAllPlayers() {
	for k, _ := range g.Players {
		g.Players[k].SetCurrentGame(g.Uid)
		g.Players[k].SendMessageOfType("team", []byte(strconv.Itoa(k)))
	}
}

func (g *GameStat) GetMapJson() []byte {
	s := g.Map.MapJson
	re := regexp.MustCompile(`\r?\n`)
	s = re.ReplaceAllString(s, " ")
	return []byte(s)
}

func (g *GameStat) GetUnitJson() []byte {
	j, _ := json.Marshal(g.Armies)
	return j
}

func (g *GameStat) PlayerDisconnect(dc *gameclient.GameClient) {
	//logic to clean up the game here, but for the time being w/e
	g.SendMessageToAllPlayers("announce", []byte("Player "+strconv.Itoa(g.GetPlayerIndex(dc))+" has left!"))
	for k, _ := range g.Players {
		if g.Players[k].IsStillConnected() && g.Players[k] != dc {
			g.EndGame(g.Players[k])
			return
		}
	}
}

func (g *GameStat) EndGame(winner *gameclient.GameClient) {
	for k, _ := range g.Players {
		if g.Players[k].IsStillConnected() {
			g.Players[k].ResetCurrentGame()
		}
	}
	ret := map[string]interface{}{"winner": strconv.Itoa(g.GetPlayerIndex(winner))}
	str, err := json.Marshal(ret)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	g.SendMessageToAllPlayers("game_end", str)
	g.SendMessageToAllPlayers("announce", []byte("Game has ended!"))
	g.Status = "ENDED"
}

func (g *GameStat) GetWeapon(wId string, owner int) *units.Weapon {
	for s := range g.Armies[owner].Squads {
		for _, u := range g.Armies[owner].Squads[s].Grunts {
			for _, w := range u.Attributes.Weapons {
				if w.Uid.String() == wId {
					return &w
				}
			}
		}
		for _, w := range g.Armies[owner].Squads[s].Leader.Attributes.Weapons {
			if w.Uid.String() == wId {
				return &w
			}
		}
	}
	return nil
}

func (g *GameStat) UnitHasSightTo(u *units.Unit, t *units.Unit) bool {
	path := los.MakeLine(u.X, u.Y, t.X, t.Y)
	for _, c := range path {
		tile := g.GetTile(c[0], c[1])
		if tile.X == u.X && tile.Y == u.Y {
			continue
		}

		if tile.X == t.X && tile.Y == t.Y {
			continue
		}

		if tile.BlocksVision || tile.Unit != nil {
			return false
		}
	}
	return true
}

func (g *GameStat) GetUnit(wId string, owner int) *units.Unit {
	for s := range g.Armies[owner].Squads {
		for u := range g.Armies[owner].Squads[s].Grunts {
			if g.Armies[owner].Squads[s].Grunts[u].Uid.String() == wId {
				return &g.Armies[owner].Squads[s].Grunts[u]
			}
		}
		if g.Armies[owner].Squads[s].Leader.Uid.String() == wId {
			return &g.Armies[owner].Squads[s].Leader
		}
	}
	return nil
}

func (g *GameStat) GetUnitGlobal(wId string) *units.Unit {
	for a := range g.Armies {
		for s := range g.Armies[a].Squads {
			for u := range g.Armies[a].Squads[s].Grunts {
				if g.Armies[a].Squads[s].Grunts[u].Uid.String() == wId {
					return &g.Armies[a].Squads[s].Grunts[u]
				}
			}
			if g.Armies[a].Squads[s].Leader.Uid.String() == wId {
				return &g.Armies[a].Squads[s].Leader
			}
		}
	}
	return nil
}

func (g *GameStat) UseWeapon(u *units.Unit, target *units.Unit, w *units.Weapon, owner int) (bool, int, int) {
	used := false
	damage := -1
	roll := 0

	if w.UsesRemaining <= 0 {
		return false, -1, 0
	}

	if g.UnitActionCounts[u.Uid.String()] <= 0 {
		return false, -1, 0
	}

	if w.NoAttack != true {
		if g.UnitCombatCounts[u.Uid.String()] <= 0 {
			return false, -1, 0
		}
	}

	if gamemap.DistanceBetweenTiles(u.X, u.Y, target.X, target.Y) > w.Rng || gamemap.DistanceBetweenTiles(u.X, u.Y, target.X, target.Y) < w.MinRng || !g.UnitHasSightTo(u, target) {
		return false, -1, 0
	}

	//pre attack
	if w.Ability == true {
		//use ability here
	}

	if w.NoAttack != true {
		used, damage, roll = g.Attack(u, target, w)
	}

	//post attack
	if w.Ability == true {
		//use ability here
	}

	return used, damage, roll
}

func (g *GameStat) AwardPoints(index int, points int) {
	g.Points[index] += points
	j, _ := json.Marshal(g.Points)
	g.SendMessageToAllPlayers("game_points_update", j)
}

func (g *GameStat) Attack(attacker *units.Unit, defender *units.Unit, w *units.Weapon) (bool, int, int) {
	attacked := false
	damage := -1
	r := dice.Roll(20)
	var attackModifier int
	var damageModifier int

	/*
		attack mod / damage mod (v) keyword
		matches and gives bonuses to attack and damage if the defender has keyword
		we go through the cycle twice, once for the unit itself, and once for the weapon being used

		+5/+4 v mounted

		-4/+4 v kingslayer

		0/-3 v dwarf
	*/
	rege, _ := regexp.Compile("((\\+\\d)|(\\-\\d)|(0))\\/((\\+\\d)|(\\-\\d)|(0))\\sv\\s\\w+")
	for _, v := range attacker.Attributes.Keywords {
		modStr := rege.MatchString(v)
		if modStr == true {
			mods := strings.Split(strings.Replace(strings.Split(v, "v")[0], " ", "", -1), "/")
			target := strings.Split(v, "v")[1]
			target = strings.Replace(target, " ", "", -1)

			for _, dv := range defender.Attributes.Keywords {
				if dv == target {
					am, _ := strconv.Atoi(mods[0])
					dm, _ := strconv.Atoi(mods[1])

					attackModifier += am
					damageModifier += dm

				}
			}
		}
	}
	for _, v := range w.Keywords {
		modStr := rege.MatchString(v)
		if modStr == true {
			mods := strings.Split(strings.Split(v, "v")[0], "/")
			target := strings.Split(v, "v")[1]
			target = strings.Replace(target, " ", "", -1)

			for _, dv := range defender.Attributes.Keywords {
				if dv == target {
					am, _ := strconv.Atoi(mods[0])
					dm, _ := strconv.Atoi(mods[1])

					attackModifier += am
					damageModifier += dm
				}
			}
		}
	}

	if ((r + w.Atk + attackModifier) + g.GetTile(attacker.X, attacker.Y).AttackModifier()) > (defender.Attributes.Def + g.GetTile(defender.X, defender.Y).DefenseModifier()) {
		damage = ((w.Dmg + damageModifier + g.GetTile(defender.X, defender.Y).DamageModifier()) - defender.Attributes.Amr)
		defender.Attributes.Hps = defender.Attributes.Hps - damage
		if defender.Attributes.Hps <= 0 {
			delete(g.UnitActionCounts, defender.Uid.String())
			delete(g.UnitCombatCounts, defender.Uid.String())
			g.GetTile(defender.X, defender.Y).Unit = nil
			g.AwardPoints(attacker.Team, defender.Value)
		}
	}
	g.UnitCombatCounts[attacker.Uid.String()] -= 1
	g.UnitActionCounts[attacker.Uid.String()] -= 1
	attacked = true
	return attacked, damage, r
}

func (g *GameStat) MoveUnit(unit *units.Unit, moves [][]int) bool {
	moved := false
	fmt.Println(unit)
	for k, v := range g.UnitActionCounts {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}
	fmt.Println(g.UnitActionCounts[unit.Uid.String()])
	if g.UnitActionCounts[unit.Uid.String()] <= 0 {
		fmt.Print("Failed at ")
		fmt.Println(1)
		return false
	}

	if gamemap.DistanceBetweenTiles(unit.X, unit.Y, moves[len(moves)-1][0], moves[len(moves)-1][1]) > unit.Attributes.Spd {
		return false
	}

	if len(moves) == 1 {
		return false
	}

	if len(moves)-1 > unit.Attributes.Spd {
		return false
	}

	pathDistance := 0
	for j := range moves {
		if j == 0 {
			continue
		}

		tt := g.GetTile(moves[j][0], moves[j][1])
		pathDistance += tt.MovementModifier()
	}

	if pathDistance > unit.Attributes.Spd {
		return false
	}

	startx := unit.X
	starty := unit.Y
	fmt.Println(moves)
	for _, m := range moves {
		if m[0] == startx && m[1] == starty {
			continue
		}
		ct := g.GetTile(unit.X, unit.Y)
		nt := g.GetTile(m[0], m[1])
		if nt.IsOpen() == false {
			fmt.Println("was mid move, tile not empty")
			return moved
		}
		if gamemap.DistanceBetweenTiles(ct.X, ct.Y, nt.X, nt.Y) > 1 {
			fmt.Println("was distance more than once")
			return moved
		}
		unit.SetPos(nt.X, nt.Y)
		nt.Unit = unit
		ct.Unit = nil
		moved = true
	}
	if moved == true {
		g.UnitActionCounts[unit.Uid.String()] = g.UnitActionCounts[unit.Uid.String()] - 1
	}
	return moved
}

func (g *GameStat) SpawnAllUnits() {
	for team := range g.Armies {
		for s := range g.Armies[team].Squads {
			for k := range g.Armies[team].Squads[s].Grunts {
				grunt := &g.Armies[team].Squads[s].Grunts[k]
				g.Map.SpawnUnitOnFirstAvailable(grunt, team)
			}
			g.Map.SpawnUnitOnFirstAvailable(&g.Armies[team].Squads[s].Leader, team)
		}
	}
}

func (g *GameStat) CheckForWinner() *gameclient.GameClient {
	var p1win int
	var p2win int
	for i, s := range g.Points {
		if s >= 30 {
			if i == 0 {
				p1win = s
			} else {
				p2win = s
			}
		}
	}
	if p1win > p2win {
		return g.Players[0]
	} else if p2win > p1win {
		return g.Players[1]
	} else {
		return nil
	}
}

func (g *GameStat) EndTurn(c *gameclient.GameClient) {
	if c != g.Players[g.CurrentTurn] {
		return
	}
	if g.CurrentTurn == 0 {
		g.CurrentTurn = 1
	} else {
		g.CurrentTurn = 0
	}
	g.ResetActions()
	for _, v := range g.Map.Map {
		for _, t := range v {
			if t.Objective != nil && t.Unit != nil && g.CurrentTurn == t.Unit.Team {
				g.AwardPoints(t.Unit.Team, 2)
			}
		}
	}
	w := g.CheckForWinner()
	if w != nil {
		g.EndGame(w)
	} else {
		g.SendMessageToAllPlayers("game_turn", []byte(strconv.Itoa(g.CurrentTurn)))
	}
}

func CreateGame(db *gorm.DB, p1 *gameclient.GameClient, p2 *gameclient.GameClient) *GameStat {
	if !reflect.DeepEqual(p1.CurrentGame, uuid.UUID{}) || !reflect.DeepEqual(p2.CurrentGame, uuid.UUID{}) {
		return nil
	}
	pary := []*gameclient.GameClient{p1, p2}

	aary := []army.Army{p1.Army, p2.Army}
	aary[0].PopulateArmy()
	aary[1].PopulateArmy()
	gmap := gamemap.GetCustomMap()
	gstat := GameStat{Armies: aary, Players: pary, Uid: uuid.NewV4(), Map: gmap}
	gstat.SendMessageToAllPlayers("announce", []byte("Game "+gstat.Uid.String()+" starting!"))
	gstat.Points = make([]int, 2)
	gstat.SetCurrentGameForAllPlayers()
	gstat.SendMessageToAllPlayers("map_data", gstat.GetMapJson())
	gstat.SpawnAllUnits()
	gstat.SendMessageToAllPlayers("game_start_army_data", gstat.GetUnitJson())
	gstat.ResetActions()
	gstat.SendMessageToAllPlayers("game_turn", []byte(strconv.Itoa(gstat.CurrentTurn)))

	return &gstat
}
