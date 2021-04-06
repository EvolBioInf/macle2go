package mau

import (
	"flag"
	"fmt"
	"os"
)

const (
	progStr     = "macle2go"
	verStr      = "0.1"
	defWin      = 10000
	defIt       = 10000
	defUpstr    = 2000
	defDownstr  = 2000
	defMinGenes = 10
)

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
	S  int      // seed for minimum number generator
	GG bool     // consider whole gene
	Pu int      // promoter upstream
	Pd int      // promoter downstream
	Cm string   // command
	O  string   // write GO terms and observed symbols to file o
	
	Files []string
}

func usage() {
	fmt.Printf("Usage: %s <command> [options]\n", progStr)
	fmt.Printf("Commands:\n")
	fmt.Printf("  annotate\tannotate macle output\n")
	fmt.Printf("  enrichment\tgene ontology (GO) enrichment analysis\n")
	fmt.Printf("  quantile\tquantile of complexity distribution\n")
	fmt.Printf("  version\tprint version\n")
	os.Exit(2)
}

func quantileUsage() {
	fmt.Printf("Usage: %s quantile options\n", progStr)
	fmt.Printf("Example: %s quantile -l 2937655681 -g 0.408679 -w 20000 -p 0.05\n", progStr)
	fmt.Printf("Options:\n")
	fmt.Printf("\t-g <NUM> gc-content\n")
	fmt.Printf("\t-l <NUM> genome length\n")
	fmt.Printf("\t-w <NUM> window length\n")
	fmt.Printf("\t-p <NUM> probability\n")
	os.Exit(2)
}

func annotateUsage() {
	fmt.Printf("Usage: %s annotate options [inputFiles]\n", progStr)
	fmt.Printf("Example %s annotate -r hsRefGene.txt -c 0.9968 -w 20000 hs_20k.mac\n", progStr)
	fmt.Printf("Options:\n")
	fmt.Printf("\t-r <FILE> refGene file, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz\n")
	fmt.Printf("\t-w <NUM> window length\n")
	fmt.Printf("\t-c <NUM>  minimum complexity\n")
	fmt.Printf("\t[-C <NUM> maximum complexity; default: no upper limit]\n")
	fmt.Printf("\t[-I <NUM> iterations; default: %d]\n", defIt)
	fmt.Printf("\t[-s <NUM> seed for random number generator; default: system-generated]\n")
	fmt.Printf("\t[-u <NUM> upstream promoter region; default: %d]\n", defUpstr)
	fmt.Printf("\t[-d <NUM> downstream promoter region; default: %d]\n", defDownstr)
	fmt.Printf("\t[-G consider whole genes; default: promoter]\n")
	os.Exit(2)
}

func enrichmentUsage() {
	fmt.Printf("Usage: %s enrichment options [inputFiles]\n", progStr)
	fmt.Printf("Example: %s enrichment -i Homo_sapiens.gene_info -g gene2go -r hsRefGene.txt -c 0.9968 -w 20000 hs_20k.mac\n", progStr)
	fmt.Printf("Options:\n")
	fmt.Printf("\t-r <FILE> refGene file, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz\n")
	fmt.Printf("\t-i <FILE> gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz\n")
	fmt.Printf("\t-g <FILE> gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz\n")
	fmt.Printf("\t-c <NUM>  minimum complexity\n")
	fmt.Printf("\t-w <NUM>  window length\n")
	fmt.Printf("\t[-C <NUM>  maximum complexity; default: no upper limit]\n")
	fmt.Printf("\t[-I <NUM> iterations; default: %d]\n", defIt)
	fmt.Printf("\t[-m <NUM>  minimum number of genes per GO-category; default: %d]\n", defMinGenes)
	fmt.Printf("\t[-s <NUM> seed for random number generator; default: system-generated]\n")
	fmt.Printf("\t[-u <NUM> upstream promoter region; default: %d]\n", defUpstr)
	fmt.Printf("\t[-d <NUM> downstream promoter region; default: %d]\n", defDownstr)
	fmt.Printf("\t[-o <FILE> print GO-terms and corresponding genes to file]\n")
	fmt.Printf("\t[-G analyze whole genes; default: promoter]\n")
	os.Exit(2)
}

func version() {
	fmt.Printf("%s %s\n", progStr, verStr)
	fmt.Printf("Written by Bernhard Haubold\n")
	fmt.Printf("Distributed under the GNU General Public License\n")
	fmt.Printf("Please send bug reports to haubold@evolbio.mpg.de.\n")
	os.Exit(2)
}

func GetArgs() Args {
	var a Args

	qc := flag.NewFlagSet("quantile",   flag.ExitOnError)
	ac := flag.NewFlagSet("annotate",   flag.ExitOnError)
	ec := flag.NewFlagSet("enrichment", flag.ExitOnError) 

	// Flags for quantile
	qc.Float64Var(&a.Qg, "g", 0, "gc-content")
	qc.Float64Var(&a.Ql, "l", 0, "genome length")
	qc.Float64Var(&a.Qw, "w", 0, "window length")
	qc.Float64Var(&a.Qp, "p", 0, "probability")
	
	// Flags for annotate
	ac.StringVar( &a.R,  "r", "", "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	ac.Float64Var(&a.C,  "c", 0,  "minimum complexity")
	ac.Float64Var(&a.CC, "C", 2,  "maximum complexity")
	ac.IntVar(    &a.II, "I", defIt,  "iterations")
	ac.IntVar(    &a.W, "w", 0,  "window length")
	ac.IntVar(    &a.S, "s", 0,  "seed for random number generator")
	ac.IntVar(    &a.Pu, "u", defUpstr, "upstream promoter region")
	ac.IntVar(    &a.Pd, "d", defDownstr, "downstream promoter region")
	ac.BoolVar(   &a.GG, "G", false, "analyze whole genes; default: promoters")

	// Flags for enrichment
	ec.StringVar( &a.R,  "r", "",    "refGene.txt, e. g. http://hgdownload.soe.ucsc.edu/goldenPath/hg38/database/refGene.txt.gz")
	ec.StringVar( &a.I,  "i", "",    "gene-info file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz")
	ec.StringVar( &a.G,  "g", "",    "gene2go file, e. g. ftp://ftp.ncbi.nih.gov/gene/DATA/gene2go.gz")
	ec.Float64Var(&a.C,  "c", 0,     "minimum complexity")
	ec.IntVar(    &a.W,  "w", 0,     "window length")
	ec.Float64Var(&a.CC, "C", 2,     "maximum complexity")
	ec.IntVar(    &a.II, "I", defIt, "iterations when computing P-values")
	ec.IntVar(    &a.M,  "m", defMinGenes,    "minimum number of genes per GO-category printed")
	ec.IntVar(    &a.S, "s", 0,  "seed for random number generator")
	ec.IntVar(    &a.Pu, "u", defUpstr, "upstream promoter region")
	ec.IntVar(    &a.Pd, "d", defDownstr, "downstream promoter region")
	ec.StringVar( &a.O, "o", "",     "file for GO/sym information")
	ec.BoolVar(   &a.GG, "G", false, "analyze whole genes; default: promoters")

	if len(os.Args) == 1 {
		usage()
	}
	flag.Usage = usage
	flag.Parse()
	switch os.Args[1] {
	case "quantile":
		qc.Usage = quantileUsage
		qc.Parse(os.Args[2:])
		a.Cm = "quantile"
	case "annotate":
		ac.Usage = annotateUsage
		ac.Parse(os.Args[2:])
		a.Files = ac.Args()
		a.Cm = "annotate"
	case "enrichment":
		ec.Usage = enrichmentUsage
		ec.Parse(os.Args[2:])
		a.Files = ec.Args()
		a.Cm = "enrichment"
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

