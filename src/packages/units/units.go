package units

import (
	"encoding/json"
	"fmt"

	"github.com/adnissen/go.uuid"
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
	Value       int
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
	Keywords  []string
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
	Keywords      []string
}

var UnitList []Unit
var SquadList []Squad
var WeaponList []Weapon

func (u *Unit) SetPos(x int, y int) {
	u.X = x
	u.Y = y
}

func LoadWeapons() {

	AddWeapon(Weapon{Id: 0, DisplayName: "Longbow", Key: "longbow", MinRng: 2, Rng: 5, Atk: 10, Dmg: 4, Uses: 15, UsesRemaining: 15})
	AddWeapon(Weapon{Id: 1, DisplayName: "Hunter's Bow", Key: "longbow", MinRng: 2, Rng: 3, Atk: 11, Dmg: 2, Uses: 99, UsesRemaining: 99})
	AddWeapon(Weapon{Id: 2, DisplayName: "Short Sword", Key: "short_sword", Rng: 1, Atk: 7, Dmg: 2, Uses: 99, UsesRemaining: 99})
	AddWeapon(Weapon{Id: 3, DisplayName: "Lance", Key: "lance", Rng: 1, Atk: 13, Dmg: 5, Uses: 5, UsesRemaining: 5})
	AddWeapon(Weapon{Id: 4, DisplayName: "Spear", Key: "spear", Rng: 1, Atk: 10, Dmg: 4, Uses: 99, UsesRemaining: 99})
	AddWeapon(Weapon{Id: 5, DisplayName: "Long Sword", Key: "longsword", Rng: 1, Atk: 10, Dmg: 3, Uses: 99, UsesRemaining: 99})
	AddWeapon(Weapon{Id: 6, DisplayName: "Goblin Dagger", Key: "dagger", Rng: 1, Atk: 6, Dmg: 5, Uses: 99, UsesRemaining: 99})
	AddWeapon(Weapon{Id: 7, DisplayName: "Goblin Shortbow", Key: "shortbow", Rng: 2, Atk: 8, Dmg: 5, Uses: 99, UsesRemaining: 99})
}

func LoadUnits() {
	//BE CAREFUL ABOUT CHANGING THIS!!
	AddUnit(Unit{Id: 0, Key: "stout_guard", DisplayName: "Bromuk", Description: "", Value: 5, Attributes: UnitAttributes{Def: 18, Amr: 2, Spd: 3, Hps: 15, WeaponIds: []int{5}}, CardText: ""})
	AddUnit(Unit{Id: 1, Key: "stout_guard", DisplayName: "Stout Guard", Description: "", Value: 3, Attributes: UnitAttributes{Def: 18, Amr: 0, Spd: 3, Hps: 13, WeaponIds: []int{5}}})

	AddUnit(Unit{Id: 2, Key: "archerguy", DisplayName: "Corath", Description: "", Value: 5, Attributes: UnitAttributes{Def: 15, Amr: 2, Spd: 4, Hps: 15, WeaponIds: []int{0, 2}}})
	AddUnit(Unit{Id: 3, Key: "archerguy", DisplayName: "Mysterious Archer", Value: 4, Description: "", Attributes: UnitAttributes{Def: 15, Amr: 0, Spd: 3, Hps: 10, WeaponIds: []int{0, 2}}})

	AddUnit(Unit{Id: 4, Key: "spearman", DisplayName: "Wild Spearman", Value: 4, Description: "", Attributes: UnitAttributes{Keywords: []string{"+2/+3 v mounted"}, Def: 16, Amr: 0, Spd: 4, Hps: 14, WeaponIds: []int{4}}})

	AddUnit(Unit{Id: 5, Key: "mounted_knight", DisplayName: "Thunder Cavalry", Value: 6, Attributes: UnitAttributes{Keywords: []string{"mounted"}, Def: 15, Amr: 1, Spd: 5, Hps: 12, WeaponIds: []int{2, 3}}})
	AddUnit(Unit{Id: 6, Key: "wild_archer", DisplayName: "Wild Hunter", Value: 4, Description: "", Attributes: UnitAttributes{Def: 14, Amr: 0, Spd: 4, Hps: 11, WeaponIds: []int{1}}})

	AddUnit(Unit{Id: 7, Key: "goblin_dagger", DisplayName: "Goblin Striker", Description: "", Value: 4, Attributes: UnitAttributes{Def: 12, Amr: 0, Spd: 4, Hps: 10, WeaponIds: []int{6}}})
	AddUnit(Unit{Id: 8, Key: "goblin_archer", DisplayName: "Goblin Bowshooter", Description: "", Value: 4, Attributes: UnitAttributes{Def: 12, Amr: 0, Spd: 4, Hps: 10, WeaponIds: []int{7}}})
}

func LoadSquads() {
	AddSquad(Squad{Id: 0, Cost: 13, DisplayName: "Bromuk's Defenders", LeaderId: 0, GruntIds: []int{1, 1, 1}, Factions: []int{0, 1, 2}})
	AddSquad(Squad{Id: 1, Cost: 10, DisplayName: "Corath's Rangers", LeaderId: 2, GruntIds: []int{3, 3}, Factions: []int{0, 1, 2}})
	AddSquad(Squad{Id: 2, Cost: 12, DisplayName: "Renegade Hunters", LeaderId: 6, GruntIds: []int{4, 4}, Factions: []int{0, 1, 2}})
	AddSquad(Squad{Id: 3, Cost: 12, DisplayName: "The Thunder Cavalry", LeaderId: 5, GruntIds: []int{5, 5}, Factions: []int{0, 1, 2}})
	AddSquad(Squad{Id: 4, Cost: 12, DisplayName: "Goblin Strike Team", LeaderId: 8, GruntIds: []int{7, 7}, Factions: []int{0, 1, 2}})
}

func AddWeapon(w Weapon) {
	WeaponList = append(WeaponList, w)
}

func AddUnit(u Unit) {
	UnitList = append(UnitList, u)
}

func AddSquad(s Squad) {
	SquadList = append(SquadList, s)
}

func GetUnit(id int) Unit {
	var ret Unit
	for _, u := range UnitList {
		if u.Id == id {
			ret = u
		}
	}
	return ret
}

func GetWeapon(id int) Weapon {
	var ret Weapon
	for _, w := range WeaponList {
		if w.Id == id {
			ret = w
		}
	}
	return ret
}

func CreateUnit(id int) Unit {
	ret := GetUnit(id)
	ret.Uid = uuid.NewV4()
	for _, w := range ret.Attributes.WeaponIds {
		nw := GetWeapon(w)
		nw.Uid = uuid.NewV4()
		ret.Attributes.Weapons = append(ret.Attributes.Weapons, nw)
	}
	return ret
}

func GetSquad(id int) Squad {
	var ret Squad
	for _, s := range SquadList {
		if s.Id == id {
			ret = s
		}
	}
	return ret
}

func CreateSquad(id int) Squad {
	ret := GetSquad(id)
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
