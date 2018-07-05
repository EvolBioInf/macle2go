package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"runtime"
	"math"
	"github.com/evolbioinf/macle2go/mau"
)

func summarizeRes(res []mau.Result) mau.Result {
	var re mau.Result
	
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

func runAnalysis(cmplx []mau.Interval, symGO mau.SymGO, args mau.Args) {
	var str []string
	var res []mau.Result
	var s int

	co, numWin, numIv, obsNumSym := mau.MergeCmplx(cmplx, args)
	c := 0
	step := int(math.Ceil(float64(args.II) / float64(runtime.NumCPU())))
	ch := make(chan bool)
	if args.Cm == "annotate" {
		for i := 0; i < args.II; i += step {
			c++
			if i + step < args.II {
				s = step
			} else {
				s = args.II - i
			}
			go func(s int) {
				r := mau.GeneEnr(cmplx, numWin, obsNumSym, s, args.S)
				res = append(res, r)
				ch <- true
			}(s)
		}
		// Make sure all threads have returned
		for i := 0; i < c; i++ {
			<-ch
		}
		r := summarizeRes(res)
		mau.PrintIntervalSym(co, numWin, numIv, obsNumSym, r.E, r.P)
	}
	if args.Cm == "enrichment" {
		uniqSym := make(map[string]bool)
		for _, iv := range co {
			for k, _ := range iv.Sym { uniqSym[k] = true }
		}
		obsGOcount := mau.GOcount(uniqSym, symGO)
		expGOcount, pVal := mau.FuncEnr(cmplx, symGO, obsNumSym, obsGOcount, args)
		for k, _ := range obsGOcount {
			str = append(str, k)
		}
		sort.Strings(str)
		gd := mau.GOdescr()
		gc := mau.GOcat()
		for _, g := range str {
			d := gd[g]
			d = strings.Replace(d, " ", "_", -1)
			if obsGOcount[g] >= args.M {
				c := gc[g]
				o := obsGOcount[g]
				e := expGOcount[g]
				fmt.Printf("%s\t%d\t%.2f\t%.2f\t%.2e\t%s\t%s\n", g, o, e, float64(o)/e, pVal[g], c, d)
			}
		}
	}
}

func progname(str string) string {
	l := strings.LastIndex(str, "/")
	return str[l+1:]
}

func printCmplx(cmplx []mau.Interval) {
	for _, c := range cmplx {
		if len(c.Sym) > 0 { fmt.Println(c) }
	}
}

func main() {
	var symGO    mau.SymGO
	var args mau.Args
	var cmplx []mau.Interval
	
	args = mau.GetArgs()
	if args.Cm == "quantile" {
		mau.Quantile(args)
		os.Exit(0)
	}
	// Get symbol/GO map from info-gene and gene2go files?
	if args.Cm == "enrichment" {
		symGO = mau.GetSymGO(args.I, args.G)
	}
	// Iterate over macle output files containing complexity data.
	if len(args.Files) == 0 {
		cmplx = mau.Cmplx("stdin", args)
		runAnalysis(cmplx, symGO, args)
	} else {
		for _, f := range args.Files {
			cmplx = mau.Cmplx(f, args)
			runAnalysis(cmplx, symGO, args)
		}
	}
}
