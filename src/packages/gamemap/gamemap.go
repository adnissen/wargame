package gamemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

type Tile struct {
	TileType   string
	Walkable   bool
	Spawn      bool
	Castle     bool
	Capturable bool
	Owner      int
	X          int
	Y          int
}

type Map struct {
	Map [][]Tile
}

var MapList = make([]Map, 1)

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

func GetMap() Map {
	return MapList[0]
}

func DistanceBetweenTiles(x1 int, y1 int, x2 int, y2 int) int {
	return int(math.Abs(float64(x1)-float64(x2)) + math.Abs(float64(y1)-float64(y2)))
}
