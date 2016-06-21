package army

import (
	"encoding/json"

	"github.com/adnissen/wargame/src/packages/units"
)

type Army struct {
	Squads []units.Squad
}

func (a Army) ArmyInformation() []byte {
	j, _ := json.Marshal(a)
	return j
}
