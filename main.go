package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/evolbioinf/macle2go/mau"
)

func runAnalysis(chrCmplx mau.ChrCmplx, chrGene mau.ChrGene, symGO mau.SymGO, args mau.Args) {
	var chrIv mau.ChrInterval
	var ivSym mau.IntervalSym
	var uniqSym    map[string]bool
	var goDescr    map[string]string
	var obsGOcount map[string]int
	var expGOcount map[string]float64
	var pVal       map[string]float64
	var e, p       float64

	fmt.Fprintf(os.Stderr, "\rGet intervals with %f <= complexity %f...", args.C, args.CC);
	chrIv = mau.GetChrInterval(chrCmplx, args)
	fmt.Fprintf(os.Stderr, "\rGet genes per interval...")
	ivSym = mau.GetIntervalSym(chrIv, chrGene)
	fmt.Fprintf(os.Stderr, "\rMake genes unique...")
	uniqSym = mau.UniqSym(ivSym)
	if args.Cm == "annotate" {
		l := len(uniqSym)
		if args.II != 0 {
			e, p = mau.GeneEnr(chrCmplx, chrGene, l, args)
		}
		mau.PrintIntervalSym(ivSym, l, e, p)
		return
	}
	goDescr = mau.GetGOdescr()
	obsGOcount = mau.GOcount(uniqSym, symGO)
	mau.Enrichment(chrCmplx, chrGene, symGO, obsGOcount, args)
	expGOcount = mau.GetExpGOcount()
	pVal = mau.GetPval()
	for g, c := range obsGOcount {
		if c >= args.M {
			e := expGOcount[g]
			f := float64(c) / e
			fmt.Printf("%s\t%d\t%.2e\t%.2f\t", g, c, e, f)
			if pVal[g] < 0 {
				fmt.Printf("<%.2e", 1.0 / float64(args.II))
			} else {
				fmt.Printf("%.2e", pVal[g])
			}
			fmt.Printf("\t%s\n", goDescr[g])
		} 
	}
}

func progname(str string) string {
	l := strings.LastIndex(str, "/")
	return str[l+1:]
}

func main() {
	var chrCmplx mau.ChrCmplx
	var chrGene  mau.ChrGene
	var symGO    mau.SymGO
	var args mau.Args
	
	version := "0.1"
	
	args = mau.GetArgs(progname(os.Args[0]), version)
	if args.Cm == "quantile" {
		mau.Quantile(args)
		os.Exit(0)
	}
	// Get chromosome/gene map from refGene file.
	chrGene = mau.GetChrGene(args.R)
	// Get symbol/GO map from info-gene and gene2go files?
	if args.Cm == "enrichment" {
		fmt.Fprintf(os.Stderr, "\rReading gene & GO data...")
		symGO = mau.GetSymGO(args.I, args.G)
	}
	// Iterate over macle output files containing complexity data.
	if len(args.Files) == 0 {
		fmt.Fprintf(os.Stderr, "\rReading macle data...")
		chrCmplx = mau.GetChrCmplx("stdin")
		runAnalysis(chrCmplx, chrGene, symGO, args)
	} else {
		for _, f := range args.Files {
			fmt.Fprintf(os.Stderr, "\rReading macle data...")
			chrCmplx = mau.GetChrCmplx(f)
			runAnalysis(chrCmplx, chrGene, symGO, args)
		}
	}
}
