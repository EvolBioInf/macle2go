package mau

import (
	"flag"
	"fmt"
	"os"
)

var progStr, verStr string

type Args struct {
	// Arguments for quantile
	Qg float64  // GC content
	Ql float64  // genome length
	Qw float64  // window length
	Qp float64  // probability

	// Arguments for annotate and enrichment
	R  string   // refGene file
	I  string   // gene info file
	G  string   // gene2go file
	II int      // iterations
	W  int      // window length
	C  float64  // minimum complexity
	CC float64  // maximum complexity
	M  int      // minimum number of genes per GO-category

	Cm string   // command
	
	Files []string
}

func usage() {
	fmt.Printf("Usage: %s <command> [options]\n", progStr)
	fmt.Println("Commands:")
	fmt.Println("  annotate\tannotate macle output")
	fmt.Println("  enrichment\tgene ontology (GO) enrichment analysis")
	fmt.Println("  quantile\tquantile of complexity distribution")
	fmt.Println("\n  version\tprint version")
	os.Exit(2)
}

func annotateUsage() {
	fmt.Printf("Usage: %s annotate options [inputFiles]\n", progStr)
	fmt.Println("Options:")
	fmt.Println("\t-r <FILE> refGene file, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	fmt.Println("\t-w <NUM> window length")
	fmt.Println("\t-c <NUM>  minimum complexity")
	fmt.Println("\t[-C <NUM> maximum complexity; default: no upper limit]")
	fmt.Printf("\t[-I <NUM> iterations; default: %d]\n", 10000)
	os.Exit(2)
}

func quantileUsage() {
	fmt.Printf("Usage: %s quantile options\n", progStr)
	fmt.Printf("Options:\n")
	fmt.Printf("\t-g <NUM> gc-content\n")
	fmt.Printf("\t-l <NUM> genome length\n")
	fmt.Printf("\t-w <NUM> window length\n")
	fmt.Printf("\t-p <NUM> probability\n")
}

func enrichmentUsage() {
	fmt.Printf("Usage: %s enrichment options [inputFiles]\n", progStr)
	fmt.Println("Options:")
	fmt.Println("\t-r <FILE> refGene file, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	fmt.Println("\t-i <FILE> gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz")
	fmt.Println("\t-g <FILE> gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz")
	fmt.Println("\t-c <NUM>  minimum complexity")
	fmt.Println("\t-w <NUM>  window length")
	fmt.Println("\t[-C <NUM>  maximum complexity; default: no upper limit]")
	fmt.Printf("\t[-I <NUM> iterations; default: %d]\n", 10000)
	fmt.Printf("\t[-m <NUM>  minimum number of genes per GO-category; default: %d]\n", 10)
}

func version() {
	fmt.Printf("%s %s\n", progStr, verStr)
	fmt.Printf("Written by Bernhard Haubold\n")
	fmt.Printf("Distributed under the GNU General Public License\n")
	fmt.Printf("Please send bug reports to haubold@evolbio.mpg.de.\n")
	os.Exit(2)
}

func GetArgs(prog, vers string) Args {
	var a Args

	qc := flag.NewFlagSet("quantile",   flag.ExitOnError)
	ac := flag.NewFlagSet("annotate",   flag.ExitOnError)
	ec := flag.NewFlagSet("enrichment", flag.ExitOnError) 

	progStr = prog
	verStr  = vers
	// Flags for quantile
	qc.Float64Var(&a.Qg, "g", 0, "gc-content")
	qc.Float64Var(&a.Ql, "l", 0, "genome length")
	qc.Float64Var(&a.Qw, "w", 0, "window length")
	qc.Float64Var(&a.Qp, "p", 0, "probability")
	
	// Flags for annotate
	ac.StringVar( &a.R,  "r", "", "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	ac.Float64Var(&a.C,  "c", 0,  "minimum complexity")
	ac.Float64Var(&a.CC, "C", 2,  "maximum complexity")
	ac.IntVar(    &a.II, "I", 10000,  "iterations")
	ac.IntVar(    &a.W, "w", 0,  "window length")

	// Flags for enrichment
	ec.StringVar( &a.R,  "r", "",    "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	ec.StringVar( &a.I,  "i", "",    "gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz")
	ec.StringVar( &a.G,  "g", "",    "gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz")
	ec.Float64Var(&a.C,  "c", 0,     "minimum complexity")
	ec.IntVar(    &a.W,  "w", 0,     "window length")
	ec.Float64Var(&a.CC, "C", 2,     "maximum complexity")
	ec.IntVar(    &a.II, "I", 10000, "iterations when computing P-values")
	ec.IntVar(    &a.M,  "m", 10,    "minimum number of genes per GO-category printed")


	if len(os.Args) == 1 {
		usage()
	}

	switch os.Args[1] {
	case "quantile":
		qc.Parse(os.Args[2:])
		qc.Usage = quantileUsage
		a.Cm = "quantile"
	case "annotate":
		ac.Parse(os.Args[2:])
		ac.Usage = annotateUsage
		a.Files = ac.Args()
		a.Cm = "annotate"
	case "enrichment":
		ec.Parse(os.Args[2:])
		ec.Usage = enrichmentUsage
		a.Files = ec.Args()
		a.Cm = "enrichment"
	case "-h":
		usage()
	case "-help":
		usage()
	case "--help":
		usage()
	case "version":
		version()
	default:
		fmt.Printf("%q is not a valid command.\n", os.Args[1])
		os.Exit(2)
	}
	if qc.Parsed() {
		if a.Qg == 0 || a.Ql == 0 || a.Qw == 0 || a.Qp == 0 {
			quantileUsage()
			os.Exit(2)
		}
	}
	if ac.Parsed() {
		if a.R == "" || a.C == 0  || a.W == 0 {
			annotateUsage()
			os.Exit(2)
		}
	}
	if ec.Parsed() {
		if a.R == "" || a.I == "" || a.G == "" || a.C == 0 || a.W == 0  {
			enrichmentUsage()
			os.Exit(2)
		}
	}
	
	return a
}
