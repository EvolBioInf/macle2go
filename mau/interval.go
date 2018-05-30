// Copyright (c) 2018 Bernhard Haubold. All rights reserved.
// Please address queries to haubold@evolbio.mpg.de.
// This program is provided under the GNU General Public License:
// https://www.gnu.org/licenses/gpl.html
package mau

import "fmt"

type Interval struct {
	Chr   string
	Start int
	End   int
	Cm    float64
}

type IntervalSlice []Interval
func (p IntervalSlice) Len() int           { return len(p) }
func (p IntervalSlice) Less(i, j int) bool { return p[i].Start < p[j].Start }
func (p IntervalSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
type ChrInterval map[string]IntervalSlice
type Sym         map[string]bool
type IntervalSym map[Interval]Sym


func PrintIntervalSym(ivSym IntervalSym, n int, e float64, p float64) {
	fmt.Printf("\r# O\tE\tP\n")
	fmt.Printf("%d\t%.4f\t%v\n", n, e, p)
	fmt.Printf("# Chr\tStart\tEnd\tLen\tC_M\tSym\n")
	for i, g := range ivSym {
		fmt.Printf("%s\t%d\t%d\t%v\t%.4f", i.Chr, i.Start, i.End, i.End - i.Start + 1, i.Cm)
		first := true
		for k, _ := range g {
			if first {
				fmt.Printf("\t%s", k)
				first = false
			} else {
				fmt.Printf(",%s", k)
			}
		}
		fmt.Printf("\n")
	}
}

func CountWindows(data ChrCmplx, threshold float64) int {
	var n int

	for _, v := range data {
		for _, c := range v {
			if c.Cm > threshold { n++ }
		}
	}

	return n
}

// GetIntervalSym maps genome intervals to sets of gene symbols.
func GetIntervalSym(chrIv ChrInterval, chrGene ChrGene) IntervalSym {
	ivGene := make(IntervalSym)

	for chr, iv := range chrIv {
		gene := chrGene[chr]
		for _, g := range gene {
			for _, i := range iv {
				sym := ivGene[i]
				if sym == nil {
					sym = make(map[string]bool)
					ivGene[i] = sym
				}
				if g.End >= i.Start && g.Start <= i.End { sym[g.Sym] = true }
			}
		}
	}
	return ivGene
}

func newIv(chr string, start, end int, cm float64) Interval {
	i := new(Interval)
	i.Chr = chr
	i.Start = start
	i.End = end
	i.Cm = cm 

	return *i
}


// GetChrInterval maps chromosomes to genome intervals with complexity greater than threshold.
// Overlapping intervals are merged.
func GetChrInterval(data ChrCmplx, a Args) ChrInterval {
	var start, end, n int
	var cm float64
	
	intervals := make(ChrInterval)
	w := int(a.W / 2.0)
	for chr, ma := range data {
		open := false
		for _, m := range ma {
			cs := m.Pos - w + 1
			ce := m.Pos + w
			if m.Cm >= a.C && m.Cm <= a.CC {
				if open {
					if cs <= end {     // extend
						cm += m.Cm
						n++
						end = ce
					} else {           // close
						open = false
						i := newIv(chr, start, end, cm/float64(n))
						intervals[chr] = append(intervals[chr], i)
					}
				} else {
					open = true
					start = cs
					end = ce
					cm = m.Cm
					n = 1
					if start < 1 { start = 1 }
				}
			} else {
				if open && cs > end {
					open = false
					i := newIv(chr, start, end, cm/float64(n))
					intervals[chr] = append(intervals[chr], i)
				}
			}
		}
		if open {
			i := new(Interval)
			i.Chr = chr
			i.Start = start
			i.End = end
			i.Cm = cm / float64(n)
			intervals[chr] = append(intervals[chr], *i)
		}
	}
	return intervals
}

// UniqGenes turns a map of intervals mapped to sets of gene symbols into a set of unique gene symbols.
func UniqSym(ivSym IntervalSym) map[string] bool {
	uniqSym := make(map[string]bool)

	for _, genes := range ivSym {
		for gene, _ := range genes {
			uniqSym[gene] = true
		}
	}

	return uniqSym
}
