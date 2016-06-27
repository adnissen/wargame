package gamemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/adnissen/wargame/src/packages/units"
)

type Tile struct {
	TileType   string
	Key        string
	OverlayKey string
	Walkable   bool
	Spawn      bool
	SpawnTeam  int
	Castle     bool
	Capturable bool
	Owner      int
	X          int
	Y          int
	Unit       *units.Unit
}

type Map struct {
	Map [][]Tile
}

var MapList = make([]Map, 100)

func LoadMaps() {
	tempMap, err := ioutil.ReadFile("src/maps/test1.json")
	if err != nil {
		fmt.Println(err)
	}
	ts := string(tempMap)
	MapList[0] = ImportMap(ts)
}

func ImportMap(s string) Map {
	res := Map{}
	json.Unmarshal([]byte(s), &res)
	return res
}

func InsertMap(m Map) {
	MapList[1] = m
	fmt.Println("imported map!")
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

			u.Spawned = true
			u.SetPos(m.Map[k][i].X, m.Map[k][i].Y)
			m.Map[k][i].Unit = u
			return
		}
	}
}

func GetMap() Map {
	return MapList[0]
}

func GetCustomMap() Map {
	if &MapList[1] != nil {
		return MapList[1]
	}
	return MapList[0]
}

func DistanceBetweenTiles(x1 int, y1 int, x2 int, y2 int) int {
	return int(math.Abs(float64(x1)-float64(x2)) + math.Abs(float64(y1)-float64(y2)))
}
