package units

import (
	"encoding/json"
	"fmt"
)

type Unit struct {
	Id          int
	DisplayName string
	Description string
	Attributes  UnitAttributes
	CardText    string
}

type Squad struct {
	Id          int
	DisplayName string
	Leader      Unit
	Grunts      []Unit
}

type UnitAttributes struct {
	Atk int
	Def int
	Amr int
	Spd int
	Hps int
	Rng int
	Dmg int
}

var UnitList = make([]Unit, 4)
var SquadList = make([]Squad, 2)

func LoadUnits() {
	//BE CAREFUL ABOUT CHANGING THIS!!
	UnitList[0] = Unit{Id: 0, DisplayName: "Bromuk", Description: "A fearless Leader, loved by friends and hated by the rest.", Attributes: UnitAttributes{Atk: 4, Def: 18, Amr: 2, Spd: 5, Hps: 15, Dmg: 6, Rng: 1}, CardText: "*Defensive Line*: "}
	UnitList[1] = Unit{Id: 1, DisplayName: "Stout Guard", Description: "You should be able to strike him down in one hit. If you can get past his shield.", Attributes: UnitAttributes{Atk: 0, Def: 20, Amr: 0, Spd: 3, Hps: 4, Dmg: 2, Rng: 1}}

	//
	UnitList[2] = Unit{Id: 2, DisplayName: "Corath", Description: "", Attributes: UnitAttributes{Atk: 4, Def: 14, Amr: 2, Spd: 6, Hps: 15, Dmg: 8, Rng: 5}}
	UnitList[3] = Unit{Id: 3, DisplayName: "Mysterious Archer", Description: "", Attributes: UnitAttributes{Atk: 1, Def: 14, Amr: 0, Spd: 4, Hps: 8, Dmg: 3, Rng: 4}}

}

func LoadSquads() {
	SquadList[0] = Squad{Id: 0, DisplayName: "Bromuk's Defenders", Leader: UnitList[0], Grunts: []Unit{UnitList[1], UnitList[1], UnitList[1]}}
	SquadList[1] = Squad{Id: 1, DisplayName: "Corath's Rangers", Leader: UnitList[2], Grunts: []Unit{UnitList[3], UnitList[3]}}
}

func GetUnit(id int) Unit {
	return UnitList[id]
}

func GetSquad(id int) Squad {
	return SquadList[id]
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
