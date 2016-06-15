package army

import "encoding/json"

type Army struct {
	Squads []int
}

func (a Army) ArmyInformation() []byte {
	j, _ := json.Marshal(a)
	return j
}
