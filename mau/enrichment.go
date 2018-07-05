package mau

import (
	"math/rand"
	"time"
)

type Result struct {
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

func GeneEnr(cmplx []Interval, nWin, nGene, it, seed int) Result {
	var p float64 
	var r *rand.Rand
	var res Result

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

func FuncEnr(cmplx []Interval, symGO SymGO, nGene int, obsGOcount map[string]int, args Args) (map[string]float64, map[string]float64) {
	var expGOcount = make(map[string]float64)
	var pVal       = make(map[string]float64)
	var r *rand.Rand

	r = seedRand(args.S)
	n := len(cmplx)

	for i := 0; i < args.II; i++ {
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
	for g := range obsGOcount {
		expGOcount[g] /= float64(args.II)
		if pVal[g] > 0 {
			pVal[g] /= float64(args.II)
		} else {
			pVal[g] = -1.0 / float64(args.II)
		}
	}
	return expGOcount, pVal
}
