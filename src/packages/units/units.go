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
	Factions    []int
}

type UnitAttributes struct {
	Def       int
	Amr       int
	Spd       int
	Hps       int
	Abilities []string
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
	MinRng        int
	Dmg           int
	Uses          int
	UsesRemaining int
	NoAttack      bool
	Ability       bool
	AbilityName   string
}

var UnitList = make([]Unit, 5)
var SquadList = make([]Squad, 2)
var WeaponList = make([]Weapon, 3)

func (u *Unit) SetPos(x int, y int) {
	u.X = x
	u.Y = y
}

func LoadWeapons() {
	WeaponList[0] = Weapon{Id: 0, DisplayName: "Longbow", Key: "longbow", Rng: 5, Atk: 13, Dmg: 10, Uses: 99, UsesRemaining: 99}
	WeaponList[1] = Weapon{Id: 1, DisplayName: "Longbow 2", Key: "longbow", Rng: 5, Atk: 10, Dmg: 3, Uses: 99, UsesRemaining: 99}
	WeaponList[2] = Weapon{Id: 2, DisplayName: "Short Sword", Key: "short_sword", Rng: 1, Atk: 5, Dmg: 10, Uses: 99, UsesRemaining: 99}
}

func LoadUnits() {
	//BE CAREFUL ABOUT CHANGING THIS!!
	UnitList[0] = Unit{Id: 0, Key: "stout_guard", DisplayName: "Bromuk", Description: "", Attributes: UnitAttributes{Def: 18, Amr: 2, Spd: 5, Hps: 15, WeaponIds: []int{2}}, CardText: "*Defensive Line*: "}
	UnitList[4] = Unit{Id: 4, Key: "stout_guard", DisplayName: "Man-At-Arms", Description: "", Attributes: UnitAttributes{Def: 15, Amr: 0, Spd: 5, Hps: 10, WeaponIds: []int{2}}}
	UnitList[1] = Unit{Id: 1, Key: "stout_guard", DisplayName: "Stout Guard", Description: "", Attributes: UnitAttributes{Def: 20, Amr: 0, Spd: 5, Hps: 4, WeaponIds: []int{2}}}

	UnitList[2] = Unit{Id: 2, Key: "archerguy", DisplayName: "Corath", Description: "", Attributes: UnitAttributes{Def: 15, Amr: 2, Spd: 6, Hps: 15, WeaponIds: []int{0}}}
	UnitList[3] = Unit{Id: 3, Key: "archerguy", DisplayName: "Mysterious Archer", Description: "", Attributes: UnitAttributes{Def: 15, Amr: 0, Spd: 6, Hps: 10, WeaponIds: []int{0}}}
}

func LoadSquads() {
	SquadList[0] = Squad{Id: 0, Cost: 15, DisplayName: "Bromuk's Defenders", LeaderId: 0, GruntIds: []int{4, 4, 4}, Factions: []int{0, 1}}
	SquadList[1] = Squad{Id: 1, Cost: 20, DisplayName: "Corath's Rangers", LeaderId: 2, GruntIds: []int{3, 3}, Factions: []int{0, 1}}
}

func GetUnit(id int) Unit {
	ret := UnitList[id]
	return ret
}

func CreateUnit(id int) Unit {
	ret := UnitList[id]
	ret.Uid = uuid.NewV4()
	for _, w := range ret.Attributes.WeaponIds {
		nw := WeaponList[w]
		fmt.Println(nw.DisplayName)
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
