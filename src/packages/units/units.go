package units

import (
	"encoding/json"
	"fmt"
)

type Unit struct {
	Id          int
	DisplayName string
	ImageUrl    string
	Description string
	Attributes  UnitAttributes
}

type Squad struct {
	Id          int
	DisplayName string
	ImageUrl    string
	Leader      int
	Grunts      []int
}

type UnitAttributes struct {
	Atk int
	Def int
	Amr int
	Spd int
}

var UnitList = make([]Unit, 2)
var SquadList = make([]Squad, 1)

func LoadUnits() {
	//BE CAREFUL ABOUT CHANGING THIS!!
	UnitList[0] = Unit{Id: 0, DisplayName: "Bromuk", ImageUrl: "https://placehold.it/32x32", Description: "A fearless Leader, loved by friends and hated by the rest.", Attributes: UnitAttributes{Atk: 5, Def: 5, Amr: 5, Spd: 5}}
	UnitList[1] = Unit{Id: 1, DisplayName: "Bromuk's Guard", ImageUrl: "https://placehold.it/32x32", Description: "Just one of the guys, I suppose.", Attributes: UnitAttributes{Atk: 6, Def: 2, Amr: 2, Spd: 4}}
}

func LoadSquads() {
	SquadList[0] = Squad{Id: 0, DisplayName: "Bastion of hope", ImageUrl: "https://placehold.it/32x32", Leader: 0, Grunts: []int{1, 1, 1}}
}

func UnitInformation() []byte {
	j, e := json.Marshal(UnitList)
	if e != nil {
		fmt.Println(e)
	}
	return j
}

func SquadInformation() []byte {
	j, e := json.Marshal(SquadList)
	if e != nil {
		fmt.Println(e)
	}
	return j
}
