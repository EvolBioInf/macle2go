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

func seedRand(args Args) *rand.Rand {
	var r *rand.Rand

	if args.S != 0 {
		r = rand.New(rand.NewSource(int64(args.S)))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	
	return r
}

func GeneEnr(cmplx []Interval, nWin, nGene int, args Args ) (float64, float64) {
	var p float64 
	var r *rand.Rand

	r = seedRand(args)
	n := len(cmplx)
	nSym := 0

	for i := 0; i < args.II; i++ {
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
	if p > 0 {
		p /= float64(args.II)
	} else {
		p = -1.0 / float64(args.II)
	}
	e := float64(nSym) / float64(args.II)

	return e, p
}

func FuncEnr(cmplx []Interval, symGO SymGO, nGene int, obsGOcount map[string]int, args Args) (map[string]float64, map[string]float64) {
	var expGOcount = make(map[string]float64)
	var pVal       = make(map[string]float64)
	var r *rand.Rand

	r = seedRand(args)
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
