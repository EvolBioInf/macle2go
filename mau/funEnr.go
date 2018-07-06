package mau

import (
	"math"
	"math/rand"
	"runtime"
)

type FunRes struct {
	E map[string]float64
	P map[string]float64
	N int
}

func sumFunRes(res []FunRes) FunRes {
	var re FunRes
	
	re.E = make(map[string]float64)
	re.P = make(map[string]float64)
	for _, r := range res {
		for k, v := range r.E {
			re.E[k] += v
			re.P[k] += r.P[k]
		}
		re.N += r.N
	}
	for k, v := range re.E {
		re.E[k] = v / float64(re.N)
		if re.P[k] > 0 {
			re.P[k] = re.P[k] / float64(re.N)
		} else {
			re.P[k] = -1. / float64(re.N)
		}
	}

	return re
}

func fEnr(cmplx []Interval, symGO SymGO, nGene int, obsGOcount map[string]int, it, seed int) (FunRes) {
	var expGOcount = make(map[string]float64)
	var pVal       = make(map[string]float64)
	var r *rand.Rand
	var fr FunRes

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
	fr.P = pVal
	fr.E = expGOcount
	fr.N = it

	return fr
}

func FunEnr(cmplx []Interval, symGO SymGO, obsNumSym int, obsGOcount map[string]int, args Args) FunRes {
	var fr  []FunRes
	var s int
	
	c := 0
	step := int(math.Ceil(float64(args.II) / float64(runtime.NumCPU())))
	ch := make(chan bool)
	for i := 0; i < args.II; i += step {
		c++
		if i + step < args.II {
			s = step
		} else {
			s = args.II - i
		}
		go func(s int) {
			r := fEnr(cmplx, symGO, obsNumSym, obsGOcount, s, args.S)
			fr = append(fr, r)
			ch <- true
		}(s)
	}
	// Make sure all threads have returned
	for i := 0; i < c; i++ {
		<-ch
	}
	r := sumFunRes(fr)
	
	return r
}
