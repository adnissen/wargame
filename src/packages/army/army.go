package army

import (
	"encoding/json"

	"github.com/adnissen/gorm"
	"github.com/adnissen/wargame/src/packages/units"
)

type Army struct {
	gorm.Model
	Squads    []units.Squad `gorm:"-"`
	SquadIds  []int         `gorm:"-"`
	IdsString string
	UserId    uint
}

func CreateArmy(db *gorm.DB, a Army) *Army {
	if err := db.Create(&a).Error; err != nil {
		return nil
	} else {
		db.Save(&a)
		return &a
	}
}

func (a Army) ArmyInformation() []byte {
	j, _ := json.Marshal(a)
	return j
}

func (a *Army) BeforeSave() (err error) {
	j, _ := json.Marshal(a.SquadIds)
	a.IdsString = string(j)
	return
}

func (a *Army) AfterFind() (err error) {
	json.Unmarshal([]byte(a.IdsString), &a.SquadIds)
	return
}

func (a *Army) PopulateArmy() {
	for s := range a.SquadIds {
		a.Squads = append(a.Squads, units.CreateSquad(s))
	}
}
