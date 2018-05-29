package mau

import (
	"os"
	"fmt"
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

func getNcmplx(cmplx map[Cmplx]bool, n int) ChrCmplx {
	var i int
	
	cc := make(ChrCmplx)
	for c, _ := range cmplx {
		i++
		if i > n { break }
		cc[c.Chr] = append(cc[c.Chr], c)
	}

	return cc
}

var expGOcount map[string]float64 = make(map[string]float64)
var pVal       map[string]float64 = make(map[string]float64)

func GetExpGOcount() map[string]float64 {
	return expGOcount
}

func GetPval() map[string]float64 {
	return pVal
}

func Enrichment(chrCmplx ChrCmplx, chrGene ChrGene, symGO SymGO, obsGOcount map[string]int, args Args) {
	var cmplx map[Cmplx]bool
	var n int
	
	cmplx, n = GetCmplx(chrCmplx, args)
	
	step := args.II / 10
	args.C = 0
	args.CC = 2
	fmt.Printf("Running Enrichment with %d iterations.\n", args.II)
	for i := 0; i < args.II; i++ {
		if (i+1) % step == 0 { fmt.Fprintf(os.Stderr, ".") }
		cm := getNcmplx(cmplx, n)
		chrIv := GetChrInterval(cm, args)
		ivSym := GetIntervalSym(chrIv, chrGene)
		uniqSym := UniqSym(ivSym)
		simGOcount := GOcount(uniqSym, symGO)
		for g, c := range simGOcount {
			expGOcount[g] += float64(c)
			if c >= obsGOcount[g] {
				pVal[g]++
			}
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
	for g := range obsGOcount {
		expGOcount[g] /= float64(args.II)
		if pVal[g] > 0 {
			pVal[g] /= float64(args.II)
		} else {
			pVal[g] = -1.0 / float64(args.II)
		}
	}
}
