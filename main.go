package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evolbioinf/macle2go/mau"
)

func runAnalysis(cmplx []mau.Interval, symGO mau.SymGO, args mau.Args) {
	// var chrIv mau.ChrInterval
	// var ivSym mau.ChrInterval
	// var uniqSym    map[string]bool
	// var goDescr    map[string]string
	// var obsGOcount map[string]int
	// var expGOcount map[string]float64
	// var pVal       map[string]float64
	// var e, p       float64
	var str []string
	//	fmt.Fprintf(os.Stderr, "\rGet intervals with %f <= complexity %f...", args.C, args.CC);
	co, w, g := mau.MergeCmplx(cmplx, args)
	if args.Cm == "annotate" {
		e, p := mau.GeneEnr(cmplx, w, g, args)
		mau.PrintIntervalSym(co, w, g, e, p)
	}
	if args.Cm == "enrichment" {
		uniqSym := make(map[string]bool)
		for _, iv := range co {
			for k, _ := range iv.Sym { uniqSym[k] = true }
		}
		obsGOcount := mau.GOcount(uniqSym, symGO)
		expGOcount, pVal := mau.FuncEnr(cmplx, symGO, g, obsGOcount, args)
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
	// for chr, iv := range(chrIv) {
	// 	fmt.Println(chr, len(iv))
	// }
	// fmt.Fprintf(os.Stderr, "\rGet genes per interval...")
	// fmt.Fprintf(os.Stderr, "\rMake genes unique...")
	// uniqSym = mau.UniqSym(chrIv)
	// if args.Cm == "annotate" {
	// 	l := len(uniqSym)
	// 	if args.II != 0 {
	// 		e, p = mau.GeneEnr(chrCmplx, l, args)
	// 	}
	// 	mau.PrintIntervalSym(ivSym, l, e, p)
	// 	return
	// }
	// goDescr = mau.GetGOdescr()
	// obsGOcount = mau.GOcount(uniqSym, symGO)
	// mau.Enrichment(chrCmplx, chrGene, symGO, obsGOcount, args)
	// expGOcount = mau.GetExpGOcount()
	// pVal = mau.GetPval()
	// for g, c := range obsGOcount {
	// 	if c >= args.M {
	// 		e := expGOcount[g]
	// 		f := float64(c) / e
	// 		fmt.Printf("%s\t%d\t%.2e\t%.2f\t%.2e", g, c, e, f, pVal[g])
	// 		fmt.Printf("\t%s\n", goDescr[g])
	// 	} 
	// }
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
	version := "0.1"
	
	args = mau.GetArgs(progname(os.Args[0]), version)
	if args.Cm == "quantile" {
		mau.Quantile(args)
		os.Exit(0)
	}
	// Get symbol/GO map from info-gene and gene2go files?
	if args.Cm == "enrichment" {
		fmt.Fprintf(os.Stderr, "Reading gene & GO data...\n")
		symGO = mau.GetSymGO(args.I, args.G)
	}
	// Iterate over macle output files containing complexity data.
	if len(args.Files) == 0 {
		fmt.Fprintf(os.Stderr, "Reading macle data from stdin...\n")
		cmplx = mau.Cmplx("stdin", args)
		runAnalysis(cmplx, symGO, args)
	} else {
		for _, f := range args.Files {
			fmt.Fprintf(os.Stderr, "Reading macle data from file...\n")
			cmplx = mau.Cmplx(f, args)
			runAnalysis(cmplx, symGO, args)
		}
	}
}
