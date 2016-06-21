package gamemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
