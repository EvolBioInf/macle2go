package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evolbioinf/macle2go/mau"
)

func runAnalysis(cmplx []mau.Interval, symGO mau.SymGO, args mau.Args) {
	var str []string
	co, numWin, numIv, obsNumSym := mau.MergeCmplx(cmplx, args)
	if args.Cm == "annotate" {
		expNumSym, p := mau.GeneEnr(cmplx, numWin, obsNumSym, args)
		mau.PrintIntervalSym(co, numWin, numIv, obsNumSym, expNumSym, p)
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
		for _, g := range str {
			d := gd[g]
			d = strings.Replace(d, " ", "_", -1)
			if obsGOcount[g] >= args.M {
				o := obsGOcount[g]
				e := expGOcount[g]
				fmt.Printf("%s\t%d\t%.2f\t%.2f\t%.2e\t%s\n", g, o, e, float64(o)/e, pVal[g], d)
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
