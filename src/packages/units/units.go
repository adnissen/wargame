package units

import (
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
)

type Unit struct {
	Id          int
	Uid         uuid.UUID
	Key         string
	DisplayName string
	Description string
	Attributes  UnitAttributes
	CardText    string
	X           int
	Y           int
	Spawned     bool
	Team        int
}

type Squad struct {
	Id          int
	Uid         uuid.UUID
	DisplayName string
	LeaderId    int
	Leader      Unit
	GruntIds    []int
	Grunts      []Unit
	Cost        int
}

type UnitAttributes struct {
	Def       int
	Amr       int
	Spd       int
	Hps       int
	WeaponIds []int
	Weapons   []Weapon
}

type Weapon struct {
	Id            int
	Uid           uuid.UUID
	DisplayName   string
	Key           string
	Atk           int
	Rng           int
	Dmg           int
	Uses          int
	UsesRemaining int
	NoAttack      bool
	Ability       bool
	AbilityName   string
}

var UnitList = make([]Unit, 4)
var SquadList = make([]Squad, 2)
var WeaponList = make([]Weapon, 2)

func (u *Unit) SetPos(x int, y int) {
	u.X = x
	u.Y = y
}

func LoadWeapons() {
	WeaponList[0] = Weapon{Id: 0, DisplayName: "Longbow", Key: "longbow", Rng: 5, Atk: 1, Dmg: 3, Uses: 20, UsesRemaining: 20}
	WeaponList[1] = Weapon{Id: 1, DisplayName: "Longbow 2", Key: "longbow", Rng: 5, Atk: 1, Dmg: 3, Uses: 20, UsesRemaining: 20}
}

func LoadUnits() {
	//BE CAREFUL ABOUT CHANGING THIS!!
	UnitList[0] = Unit{Id: 0, Key: "bromuk", DisplayName: "Bromuk", Description: "A fearless Leader, loved by friends and hated by the rest.", Attributes: UnitAttributes{Def: 18, Amr: 2, Spd: 5, Hps: 15, WeaponIds: []int{1}}, CardText: "*Defensive Line*: "}
	UnitList[1] = Unit{Id: 1, Key: "stout_guard", DisplayName: "Stout Guard", Description: "You should be able to strike him down in one hit. If you can get past his shield.", Attributes: UnitAttributes{Def: 20, Amr: 0, Spd: 3, Hps: 4, WeaponIds: []int{1}}}

	UnitList[2] = Unit{Id: 2, Key: "corath", DisplayName: "Corath", Description: "", Attributes: UnitAttributes{Def: 14, Amr: 2, Spd: 6, Hps: 15, WeaponIds: []int{1}}}
	UnitList[3] = Unit{Id: 3, Key: "mysterious_archer", DisplayName: "Mysterious Archer", Description: "", Attributes: UnitAttributes{Def: 14, Amr: 0, Spd: 4, Hps: 8, WeaponIds: []int{1}}}
}

func LoadSquads() {
	SquadList[0] = Squad{Id: 0, Cost: 5, DisplayName: "Bromuk's Defenders", LeaderId: 0, GruntIds: []int{1, 1, 1}}
	SquadList[1] = Squad{Id: 1, Cost: 5, DisplayName: "Corath's Rangers", LeaderId: 2, GruntIds: []int{3, 3, 3}}
}

func GetUnit(id int) Unit {
	ret := UnitList[id]
	return ret
}

func CreateUnit(id int) Unit {
	ret := UnitList[id]
	ret.Uid = uuid.NewV4()
	for w := range ret.Attributes.WeaponIds {
		nw := WeaponList[w]
		nw.Uid = uuid.NewV4()
		ret.Attributes.Weapons = append(ret.Attributes.Weapons, nw)
	}
	return ret
}

func GetSquad(id int) Squad {
	ret := SquadList[id]
	return ret
}

func CreateSquad(id int) Squad {
	ret := SquadList[id]
	ret.Uid = uuid.NewV4()
	for i := range ret.GruntIds {
		ret.Grunts = append(ret.Grunts, CreateUnit(ret.GruntIds[i]))
	}
	ret.Leader = CreateUnit(ret.LeaderId)
	return ret
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
