package mau

import (
	"math"
	"math/rand"
	"runtime"
)

type AnnRes struct {
	E float64
	P float64
	N int
}

func logBinom(n, k int64) float64 {
	x1, _ := math.Lgamma(float64(n + 1))
	x2, _ := math.Lgamma(float64(k + 1))
	x3, _ := math.Lgamma(float64(n - k + 1))
	x := x1 - x2 - x3
	
	return x
}

// AnnCon computes the contingency of symbol annotations
func AnnCon(cmplx []Interval, args Args) (int64, int64, int64, int64, float64)  {
	var a int64 //  hiCm /  sym
	var b int64 //  hiCm / !sym
	var c int64 // !hiCm /  sym
	var d int64 // !hiCm / !sym
	var n int64
	var co, sy bool

	for _, i := range cmplx {
		if i.Cm >= args.C {
			co = true
		} else {
			co = false
		}
		if len(i.Sym) > 0 {
			sy = true
		} else {
			sy = false
		}
		if co && sy {
			a++
		} else if co && !sy {
			b++
		} else if !co && sy {
			c++
		} else if !co && !sy {
			d++
		}
		
	}
	n = a + b + c + d
	p := logBinom(a + b, a) + logBinom(c + d, c) - logBinom(n, a + c)
	
	return a, b, c, d, math.Exp(p)
}

func sumAnnRes(res []AnnRes) AnnRes {
	var re AnnRes
	
	for _, r := range res {
		re.N += r.N
		re.P += r.P
		re.E += r.E
	}
	if re.P == 0 {
		re.P = -1. / float64(re.N)
	} else {
		re.P /= float64(re.N)
	}
	if re.E == 0 {
		re.E = -1. / float64(re.N)
	} else {
		re.E /= float64(re.N)
	}

	return re
}

func aEnr(cmplx []Interval, nWin, nGene, it, seed int) AnnRes {
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

func AnnEnr(cmplx []Interval, nWin, nGene int, args Args) AnnRes {
	var ar  []AnnRes
	
	c := 0
	step := int(math.Ceil(float64(args.II) / float64(runtime.NumCPU())))
	ch := make(chan bool)
	var s int

	for i := 0; i < args.II; i += step {
		c++
		if i + step < args.II {
			s = step
		} else {
			s = args.II - i
		}
		go func(s int) {
			r := aEnr(cmplx, nWin, nGene, s, args.S)
			ar = append(ar, r)
			ch <- true
		}(s)
	}
	// Make sure all threads have returned
	for i := 0; i < c; i++ {
		<-ch
	}
	r := sumAnnRes(ar)
	
	return r
}
