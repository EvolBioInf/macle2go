package main

import (
	"fmt"
	"os"
	"flag"
	"log"

	"github.com/evolbioinf/macle2go"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var macle    []macle2go.Macle
	var refGene  []macle2go.RefGene
	var geneInfo []macle2go.GeneInfo
	var gene2go  []macle2go.Gene2go
	var err error
	var f   *os.File
	var ref  = flag.String("r", "", "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	var info = flag.String("i", "", "gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz")
	var gg   = flag.String("g", "", "gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz")
	var t    = flag.Float64("t", 0, "complexity threshold")
	var w    = flag.Int("w", 0, "window length")

	flag.Parse()
	if *ref == "" || *info == "" || *gg == "" || *t == 0 || *w == 0 {
		fmt.Println("Please enter all the options for running macle2go.")
		flag.Usage()
		os.Exit(2)
	}
	// Read the three data files refGene, geneInfo, and gene2go
	f, err = os.Open(*ref)
	check(err)
	refGene = macle2go.GetRefGene(f)
	f.Close()
	fmt.Println("refGene done: ", refGene[0])
	f, err = os.Open(*info)
	check(err)
	geneInfo = macle2go.GetGeneInfo(f)
	f.Close()
	fmt.Println("geneInfo done: ", geneInfo[0])
	f, err = os.Open(*gg)
	check(err)
	gene2go = macle2go.GetGene2go(f)
	f.Close()
	fmt.Println("gene2go done: ", gene2go[0])
	// Iterate over the macle input files
	//	files := os.Args[len(flag.Args())]
	files := os.Args[11:]
	if len(files) == 0 {
		macle = macle2go.GetMacle(os.Stdin)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "macle2func: %v\n", err)
				continue
			}
			macle = macle2go.GetMacle(f)
			f.Close()
		}
	}
	n := macle2go.CountWindows(macle, *t)
	fmt.Println("Number of windows: ", n)
	intervals := macle2go.GetIntervals(macle, *t, *w)
	im := macle2go.IntervalMap(intervals)
	fmt.Println("First interval: ", intervals[0])
}
