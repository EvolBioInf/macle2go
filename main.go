package main

import (
	"fmt"
	"os"
	
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

	chrIv = mau.GetChrInterval(chrCmplx, args)
	ivSym = mau.GetIntervalSym(chrIv, chrGene)
	if args.P {
		mau.PrintIntervalSym(ivSym)
		os.Exit(0)
	}
	uniqSym = mau.UniqSym(ivSym)
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

func main() {
	var chrCmplx mau.ChrCmplx
	var chrGene  mau.ChrGene
	var symGO    mau.SymGO
	var args mau.Args
	
	args = mau.GetArgs()
	// Get chromosome/gene map from refGene file.
	chrGene = mau.GetChrGene(args.R)
	// Get symbol/GO map from info-gene and gene2go files.
	symGO = mau.GetSymGO(args.I, args.G)
	// Iterate over macle output files containing complexity data.
	if len(args.Files) == 0 {
		chrCmplx = mau.GetChrCmplx("stdin")
		runAnalysis(chrCmplx, chrGene, symGO, args)
	} else {
		for _, f := range args.Files {
			chrCmplx = mau.GetChrCmplx(f)
			runAnalysis(chrCmplx, chrGene, symGO, args)
		}
	}
}
