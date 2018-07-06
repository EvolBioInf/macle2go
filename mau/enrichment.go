package mau

import (
	"math/rand"
	"time"
)

func GOcount(uniqSym map[string]bool, symGO SymGO) map[string]int {
	var goCount map[string]int

	goCount = make(map[string]int)
	for sym, _ := range uniqSym {
		goAcc := symGO[sym]
		for g, _ := range goAcc {
			goCount[g]++
		}
	}
	
	return goCount
}

func seedRand(seed int) *rand.Rand {
	var r *rand.Rand

	if seed != 0 {
		r = rand.New(rand.NewSource(int64(seed)))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	
	return r
}
