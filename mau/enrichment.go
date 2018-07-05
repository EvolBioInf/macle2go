package mau

import (
	"math/rand"
	"time"
)

type EnrRes struct {
	E map[string]float64
	P map[string]float64
	N int
}

type AnnRes struct {
	E float64
	P float64
	N int
}

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

func GeneEnr(cmplx []Interval, nWin, nGene, it, seed int) AnnRes {
	var p float64 
	var r *rand.Rand
	var res AnnRes

	r = seedRand(seed)
	n := len(cmplx)
	nSym := 0

	for i := 0; i < it; i++ {
		sym := make(map[string]bool)
		for j := 0; j < nWin; j++ {
			x := r.Intn(n)
			for k, _ := range cmplx[x].Sym {
				sym[k] = true
			}
		}
		nSym += len(sym)
		if len(sym) >= nGene { p++ }
	}
	res.P = p
	res.E = float64(nSym)
	res.N = it
	
	return res
}

func FuncEnr(cmplx []Interval, symGO SymGO, nGene int, obsGOcount map[string]int, it, seed int) (EnrRes) {
	var expGOcount = make(map[string]float64)
	var pVal       = make(map[string]float64)
	var r *rand.Rand
	var er EnrRes

	r = seedRand(seed)
	n := len(cmplx)

	for i := 0; i < it; i++ {
		sym := make(map[string]bool)
		for ; len(sym) < nGene; {
			x := r.Intn(n)
			for k, _ := range cmplx[x].Sym { sym[k] = true }
		}
		simGOcount := GOcount(sym, symGO)
		for g, c := range simGOcount {
			expGOcount[g] += float64(c)
			if c >= obsGOcount[g] {
				pVal[g]++
			}
		}
		
	}
	er.P = pVal
	er.E = expGOcount
	er.N = it

	return er
}
