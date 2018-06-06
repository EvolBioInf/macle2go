// Copyright (c) 2018 Bernhard Haubold. All rights reserved.
// Please address queries to haubold@evolbio.mpg.de.
// This program is provided under the GNU General Public License:
// https://www.gnu.org/licenses/gpl.html
package mau

import (
	"fmt"
	"sort"
)

type Interval struct {
	Chr   string
	Start int
	End   int
	Cm    float64
	Sym   map[string]bool
}

func NewInterval(chr string, start, end int, cm float64, sym string, s map[string]bool) *Interval {
	i := new(Interval)
	i.Chr = chr
	i.Start = start
	i.End = end
	i.Cm = cm

	if sym != "" || s != nil { i.Sym = make(map[string]bool) }
	if sym != "" { i.Sym[sym] = true }
	if s != nil {
		for k, _ := range s { i.Sym[k] = true }
	}
	return i
}

type IntervalSlice []Interval
func (p IntervalSlice) Len() int           { return len(p) }
func (p IntervalSlice) Less(i, j int) bool { return p[i].Start < p[j].Start }
func (p IntervalSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
type ChrInterval map[string]IntervalSlice

// SortStringSet takes a map of strings to bools and returns a sorted slice of the keys
func SortStringSet(sm map[string]bool) []string {
	var str []string
	for s, _ := range sm {
		str = append(str, s)
	}
	sort.Strings(str)
	return str
}

func PrintIntervalSym(iv []Interval, w, o int, e float64, p float64) {
	first := true
	fmt.Printf("# W\tO\tE\tO/E\tP\n")
	fmt.Printf("# %d\t%d\t%.2f\t%.2f\t%v\n", w, o, e, float64(o)/e, p)
	fmt.Printf("# Chr\tStart\tEnd\tLen\tC_M\tSym\n")
	for _, i := range iv {
		fmt.Printf("%s\t%d\t%d\t%v\t%.4f", i.Chr, i.Start, i.End, i.End - i.Start + 1, i.Cm)
		sy := SortStringSet(i.Sym)
		first = true
		for _, s := range(sy) {  // symbols
			if first {
				fmt.Printf("\t%s", s)
				first = false
			} else {
				fmt.Printf(" %s", s)
			}
		}
		fmt.Printf("\n")

	}
}

func CountWindows(data ChrInterval, threshold float64) int {
	var n int

	for _, v := range data {
		for _, c := range v {
			if c.Cm > threshold { n++ }
		}
	}

	return n
}

// AddSym traverses the genes and adds their symbols to the intersecting intervals
func AddSym(intervals, genes ChrInterval) {
	var sym string 
	for chr, iv := range intervals {
		gene := genes[chr]
		for _, g := range gene {
			for sym = range g.Sym {}
			for _, i := range iv {
				if g.End >= i.Start && g.Start <= i.End { i.Sym[sym] = true }
			}
		}
	}
}

// func newIv(chr string, start, end int, cm float64) Interval {
// 	i := NewInterval()
// 	i.Chr = chr
// 	i.Start = start
// 	i.End = end
// 	i.Cm = cm 

// 	return *i
// }


// GetChrInterval maps chromosomes to genome intervals with complexity greater than threshold.
// Overlapping intervals are merged.
// func GetChrInterval(data ChrCmplx, a Args) ChrInterval {
// 	var start, end, n int
// 	var cm float64
// 	var iv Interval


// 	intervals := make(ChrInterval)
// 	w := int(a.W / 2.0)
// 	for chr, ma := range data {
// 		open := false
// 		for _, m := range ma {
// 			cs := m.Pos - w + 1
// 			ce := m.Pos + w
// 			sy := m.Sym
// 			if m.Cm >= a.C && m.Cm <= a.CC {
// 				if open {
// 					if cs <= end {     // extend
// 						iv.End = ce
// 						iv.Cm += m.Cm
// 						n++
// 					} else {           // close
// 						open = false
// 					}
// 				} else {                   // open
// 					if start < 1 { start = 1 }
// 					iv = newIv(chr, start, end, m.Cm, m.Sym)
// 					open = true
// 					n = 1
// 				}
// 			} else {
// 				if open && cs > end {
// 					open = false
// 					iv = newIv(chr, start, end, cm/float64(n))
// 					for s, _ := range sy { iv.Sym[s] = true }
// 					intervals[chr] = append(intervals[chr], iv)
// 				}
// 			}
// 		}
// 		if open {
// 			i := NewInterval()
// 			i.Chr = chr
// 			i.Start = start
// 			i.End = end
// 			i.Cm = cm / float64(n)
// 			intervals[chr] = append(intervals[chr], *i)
// 		}
// 	}
// 	return intervals
// }

// UniqGenes extracts the set of gene symbols contained in intervals
func UniqSym(intervals ChrInterval) map[string] bool {
	uniqSym := make(map[string]bool)

	for _, iv := range intervals {
		for _, i := range iv {
			for s, _ := range i.Sym {
				uniqSym[s] = true
			}
		}
	}

	return uniqSym
}
