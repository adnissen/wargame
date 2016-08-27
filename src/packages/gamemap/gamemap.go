package gamemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/adnissen/wargame/src/packages/units"
)

type Tile struct {
	TileType     string
	OverlayKey   string
	Walkable     bool
	Spawn        bool
	SpawnTeam    int
	Castle       bool
	Capturable   bool
	BlocksVision bool
	Objective    *Objective
	Owner        int
	X            int
	Y            int
	Unit         *units.Unit
	AttackMod    int
	DefenseMod   int
	DamageMod    int
	MovementMod  int
}

type Map struct {
	Map     [][]Tile
	MapJson string
}

type Objective struct {
	Owned bool
	Owner int
}

var MapList = make(map[string]Map)

func LoadMaps() {
	tempMap, err := ioutil.ReadFile("src/maps/test1.json")
	if err != nil {
		fmt.Println(err)
	}
	ts := string(tempMap)
	MapList["default"] = ImportMap(ts)
}

func ImportMap(s string) Map {
	SaveMap(s)
	var dat map[string]interface{}
	json.Unmarshal([]byte(s), &dat)
	var res = Map{}
	res.Map = make([][]Tile, int(dat["width"].(float64)))
	for i := range res.Map {
		res.Map[i] = make([]Tile, int(dat["height"].(float64)))
	}
	layers := dat["layers"].([]interface{})
	for _, l := range layers {
		layer := l.(map[string]interface{})
		props, hasProps := layer["properties"].(map[string]interface{})
		if hasProps {
			for k, t := range layer["data"].([]interface{}) {
				if t.(float64) != 0 {
					tile, tx, ty := res.GetTileByIndex(k)
					if val, ok := props["Walkable"]; ok {
						tile.Walkable = val.(bool)
					}
					if val, ok := props["Spawn"]; ok {
						tile.Spawn = val.(bool)
					}
					if val, ok := props["SpawnTeam"]; ok {
						tile.SpawnTeam = int(val.(float64))
					}
					if val, ok := props["TileType"]; ok {
						tile.TileType = val.(string)
					}
					if val, ok := props["BlocksVision"]; ok {
						tile.BlocksVision = val.(bool)
					}
					if val, ok := props["DamageModifier"]; ok {
						tile.DamageMod = int(val.(float64))
					}
					if val, ok := props["AttackModifier"]; ok {
						tile.AttackMod = int(val.(float64))
					}
					if val, ok := props["DefenseModifier"]; ok {
						tile.DefenseMod = int(val.(float64))
					}
					if val, ok := props["MovementModifier"]; ok {
						tile.MovementMod = int(val.(float64))
					}
					tile.X = tx
					tile.Y = ty

					if _, ok := props["Objective"]; ok {
						tile.Objective = &Objective{}
					}
				}
			}
		}
	}

	res.MapJson = s
	return res
}

func InsertMap(m Map) {
	MapList["save"] = m
	fmt.Println("imported map!")
}

func SaveMap(s string) {
	err := ioutil.WriteFile("src/maps/test1.json", []byte(s), 0644)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("saved map!")
}

func (t *Tile) DefenseModifier() int {
	return t.DefenseMod
}

func (t *Tile) AttackModifier() int {
	return t.AttackMod
}

func (t *Tile) DamageModifier() int {
	return t.DamageMod
}

func (t *Tile) MovementModifier() int {
	if t.MovementMod == 0 {
		return 1
	}
	return t.MovementMod
}

func (m *Map) GetTileByIndex(index int) (*Tile, int, int) {
	var count int
	for i := range m.Map[0] {
		for j := range m.Map {
			if count == index {
				return &m.Map[j][i], j, i
			}
			count = count + 1
		}
	}
	return nil, 0, 0
}

func (t *Tile) IsOpen() bool {
	return t.Walkable && t.Unit == nil
}

func (m *Map) SpawnUnitOnFirstAvailable(u *units.Unit, team int) {
	for k := range m.Map {
		for i := range m.Map[k] {
			//just find the first spawn points and put the units in them
			if !m.Map[k][i].Spawn {
				continue
			}
			if m.Map[k][i].SpawnTeam != team {
				continue
			}
			if !m.Map[k][i].Walkable {
				continue
			}
			if m.Map[k][i].Unit != nil {
				continue
			}

			u.Team = team
			u.Spawned = true
			u.SetPos(m.Map[k][i].X, m.Map[k][i].Y)
			m.Map[k][i].Unit = u
			return
		}
	}
}

func GetMap() Map {
	d := MapList["default"]
	return d
}

func GetCustomMap() Map {
	tempMap, err := ioutil.ReadFile("src/maps/test1.json")
	if err != nil {
		fmt.Println(err)
	}
	return ImportMap(string(tempMap))
}

func DistanceBetweenTiles(x1 int, y1 int, x2 int, y2 int) int {
	return int(math.Abs(float64(x1)-float64(x2)) + math.Abs(float64(y1)-float64(y2)))
}
