package mau

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	R  string   // name of refGene file
	I  string   // name of gene-info file
	G  string   // name of gene2go file
	C  float64  // minimum complexity
	CC float64  // maximum complexity
	W  int      // window length
	II int      // number of iterations
	P  bool     // print gene list & exit
	M  int      // minimum number of genes per GO-category printed
	Files []string
}

func GetArgs() Args {
	var a Args

	flag.StringVar(&a.R, "r", "", "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	flag.StringVar(&a.I, "i", "", "gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz")
	flag.StringVar(&a.G, "g", "", "gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz")
	flag.Float64Var(&a.C, "c", 0, "minimum complexity")
	flag.Float64Var(&a.CC, "C", 2, "maximum complexity")
	flag.IntVar(&a.W, "w", 0, "window length")
	flag.IntVar(&a.II, "I", 10000, "iterations when computing P-values")
	flag.BoolVar(&a.P, "p", false, "print interval/gene list and exit")
	flag.IntVar(&a.M, "m", 10, "minimum number of genes per GO-category printed")
	flag.Parse()

	missing := ""
	if a.R == "" { missing += " -r" }
	if a.I == "" { missing += " -i" }
	if a.G == "" { missing += " -g" }
	if a.C == 0  { missing += " -c" }
	if a.W == 0  { missing += " -w" }
	
	if missing != "" {
		fmt.Printf("The following mandatory options have not been set %s.\n", missing)
		flag.Usage()
		os.Exit(2)
	}

	a.Files = flag.Args()

	return a
}
