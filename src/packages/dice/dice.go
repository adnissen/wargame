package dice

import (
	"math/rand"
	"time"
)

func Roll(sides int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return (r1.Intn(sides-1) + 1)
}
