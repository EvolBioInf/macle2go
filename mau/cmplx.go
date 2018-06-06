package mau

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sort"
)

var refGeneData []Interval

// Cmplx generates annotated complexity data
func Cmplx(macleFile string, args Args) []Interval {
	refGeneFile := args.R
	if refGeneData == nil {
		refGeneData = refGene(refGeneFile)
	}
	m, s, c := macle(macleFile, args.W)
	annotate(m, refGeneData, args.W, s)
	ma := rmNeg(m, c)
	return ma
}

// MergeCmplx extracts and merges intervals with where args.C <= complexity <= args.CC
func MergeCmplx(cmplx []Interval, args Args) ([]Interval, int) {
	var co []Interval
	var iv *Interval
	open := false
	n := 0
	count := 0

	for _, i := range cmplx {
		if i.Cm >= args.C && i.Cm <= args.CC {
			if open {                // extend
				if i.Start <= iv.End {
					iv.End = i.End
					iv.Cm += i.Cm
					if iv.Sym == nil { iv.Sym = make(map[string]bool) }
					for k, _ := range i.Sym { iv.Sym[k] = true }
					n++
					count++
				} else {
					open = false
				}
			} else {                 // open
				iv = NewInterval(i.Chr, i.Start, i.End, i.Cm, "", i.Sym)
				open = true
				n = 1
				count++
			}
		} else {
			if open && i.Start > iv.End {               // close
				open = false
				iv.Cm /= float64(n)
				co = append(co, *iv)
			} 
		}
	}
	return co, count
}

func rmNeg(macle ChrInterval, chr []string) []Interval {
	var ma []Interval

	for _, c := range chr {
		iv := macle[c]
		for _, i := range iv {
			if i.Cm > -1 { ma = append(ma, i) }
		}
	}

	return ma
}

// annotate adds gene symbols to macle intervals
func annotate(macle ChrInterval, refGene []Interval, win, step int) {
	for _, g := range refGene {
		ma := macle[g.Chr]
		if ma == nil { continue }
		s := sort.Search(len(ma), func(i int) bool { return ma[i].End >= g.Start })
		e := sort.Search(len(ma), func(i int) bool { return ma[i].Start >= g.End })
		foundS := false
		foundE := false
		for i := s; i <= e; i++ {
			if ma[i].Start <= g.End && ma[i].End >= g.Start {
				if g.Start >= ma[i].Start && g.Start <= ma[i].End { foundS = true }
				if g.End   >= ma[i].Start && g.End   <= ma[i].End { foundE = true }
				if ma[i].Sym == nil { ma[i].Sym = make(map[string]bool) }
				for k, _ := range g.Sym {
					ma[i].Sym[k] = true
				}
			}
		}
		if foundS != true || foundE != true {
			fmt.Println("Warning macle2go: Could not find macle-interval for gene ", g)
		}
	}
}

// macle reads complexity data from a macle file
func macle(fileName string, win int) (ChrInterval, int, []string) {
	var str1, str2 []string
	const numCol = 3
	var file *os.File
	var err error
	var cs []string

	ca := make(map[string]bool)
	if fileName == "stdin" {
		file = os.Stdin
	} else {
		file, err = os.Open(fileName)
	}
	Check(err)
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	Check(err)
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Determine window length
	str2 = strings.Split(str1[0], "\t")
	p1, err := strconv.Atoi(str2[1])
	Check(err)
	// Determine step length
	str2 = strings.Split(str1[1], "\t")
	p2, err := strconv.Atoi(str2[1])
	s := p2 - p1
	// Allocate memory 
	chrIv := make(ChrInterval)
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines and comment lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		// Read chromosome, position, and complexity
		chr := str2[0]
		if ca[chr] == false {
			ca[chr] = true
			cs = append(cs, chr)
		}
		pos, err := strconv.Atoi(str2[1])
		Check(err)
		com, err := strconv.ParseFloat(str2[2], 64)
		Check(err)
		start := pos - win/2 + 1
		end := pos + win/2
		co := NewInterval(chr, start, end, com, "", nil)
		chrIv[chr] = append(chrIv[chr], *co)
	}
	file.Close()
	return chrIv, s, cs
}


// refGene reads a refGene file and returns a map of chromosomes to slices of intervals denoting genes
func refGene(fileName string) []Interval {
	var str1, str2 []string
	var chr string
	var start, end int
	var file *os.File
	var err error
	var rg []Interval

	file, err = os.Open(fileName)
	Check(err)
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	Check(err)
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Iterate over lines
	for _, s := range str1 {
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) != 16 {
			continue
		}
		chr = str2[2]
		start, err = strconv.Atoi(str2[4])
		Check(err)
		end, err = strconv.Atoi(str2[5])
		Check(err)
		gene := NewInterval(chr, start, end, 0, str2[12], nil)
		rg = append(rg, *gene)
	}
	file.Close()
	return rg
}
