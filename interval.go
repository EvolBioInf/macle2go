// Copyright (c) 2018 Bernhard Haubold. All rights reserved.
// Please address queries to haubold@evolbio.mpg.de.
// This program is provided under the GNU General Public License:
// https://www.gnu.org/licenses/gpl.html

package macle2go

type Interval struct {
	Chr   string
	Start int
	End   int
}

func CountWindows(data []Macle, threshold float64) int {
	var n int

	for _, m := range data {
		if m.Cm > threshold { n++ }
	}

	return n
}

func GetIntervals(data []Macle, threshold float64, winLen int) []Interval {
	var open bool
	var intervals []Interval
	var start, end int
	
	w := int(winLen / 2.0)
	for _, m := range data {
		if m.Cm > threshold {
			end = m.Pos + w
			if !open {
				open = true
				start = m.Pos - w
				if start < 1 { start = 1 }
				end = m.Pos + w
			} else {
				end = m.Pos + w
			}
		} else {
			if open {
				open = false
				i := new(Interval)
				i.Chr = m.Chr
				i.Start = start
				i.End = end
				intervals = append(intervals, *i)
			}
		}
	}
	return intervals
}
