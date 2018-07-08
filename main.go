package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"bufio"
	
	"github.com/evolbioinf/macle2go/mau"
)

func printGOsym(fileName string, uniqSym map[string]bool, symGO mau.SymGO, args mau.Args) {
	var str []string
	var go2sy map[string][]string

	go2sy = make(map[string][]string)
	gc := mau.CountGO(symGO)
	for s, _ := range uniqSym {
		gg := symGO[s]
		for g, _ := range gg {
			go2sy[g] = append(go2sy[g], s)
		}
	}
	for g, _ := range go2sy {
		str = append(str, g)
		sort.Strings(go2sy[g])
	}
	
	sort.Strings(str)
	f, err := os.Create(fileName)
	mau.Check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, g := range str {
		l := len(go2sy[g])
		if l >= args.M {
			fmt.Fprintf(w, "%s\t%d\t%d", g, gc[g], len(go2sy[g]))
			for _, s := range go2sy[g] {
				fmt.Fprintf(w, " %s", s)
			}
			fmt.Fprintf(w, "\n")
		}
	}
	w.Flush()
}

func runAnalysis(cmplx []mau.Interval, symGO mau.SymGO, args mau.Args) {
	var str []string

	co, numWin, numIv, obsNumSym := mau.MergeCmplx(cmplx, args)
	if args.Cm == "annotate" {
		a, b, c, d, p := mau.AnnCon(cmplx, args)
		fmt.Printf("#\tprom\t!prom\n")
		fmt.Printf("# compl\t%d\t%d\n", a, b)
		fmt.Printf("#!compl\t%d\t%d\n", c, d)
		fmt.Printf("#P: %.3e\n", p)
		r := mau.AnnEnr(cmplx, numWin, obsNumSym, args)
		mau.PrintIntervalSym(co, numWin, numIv, obsNumSym, r.E, r.P)
	}
	if args.Cm == "enrichment" {
		uniqSym := make(map[string]bool)
		for _, iv := range co {
			for k, _ := range iv.Sym { uniqSym[k] = true }
		}
		if args.O != "" {
			printGOsym(args.O, uniqSym, symGO, args)
		}
		obsGOcount := mau.GOcount(uniqSym, symGO)
		ere := mau.FunEnr(cmplx, symGO, obsNumSym, obsGOcount, args)
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
				e := ere.E[g]
				p := ere.P[g]
				fmt.Printf("%s\t%d\t%.2f\t%.2f\t%.2e\t%s\t%s\n", g, o, e, float64(o)/e, p, c, d)
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
