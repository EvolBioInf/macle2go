package mau

import (
	"strconv"
	"strings"
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
)

type Cmplx struct {
	Chr string
	Pos int
	Cm  float64
}

type ChrCmplx map[string][]Cmplx

func colCount(col []string) int {
	c := 0
	for _, s := range col {
		if len(s) > 0 { c++ }
	}
	return c
}

// GetCmplx constructs a set of all complexity values in chrCmplx.
func GetCmplx(chrCmplx ChrCmplx, a Args) (map[Cmplx]bool, int) {
	cmplx := make(map[Cmplx]bool)
	var n int

	for _, cm := range chrCmplx {
		for _, c := range cm {
			cmplx[c] = true
			if c.Cm >= a.C && c.Cm <= a.CC { n++ }
			
		}
	}
	return cmplx, n
}

// GetChrCmplx reads macle output and returns a map of chromosomes to a map of positions to complexity values.
func GetChrCmplx(fileName string) ChrCmplx {
	var str1, str2 []string
	const numCol = 3
	var file *os.File
	var err error

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
	// Allocate memory 
	chrCmplx := make(ChrCmplx)
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines and comment lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		c := colCount(str2)
		if c != numCol {
			fmt.Fprintf(os.Stderr, "Warning[GetMacle]: Expect %d columns, but find %d in %s\n", numCol, len(str2), s)
			continue
		}
		// Read chromosome, position, and compolexity
		chr := str2[0]
		pos, err := strconv.Atoi(str2[1])
		Check(err)
		com, err := strconv.ParseFloat(str2[2], 64)
		Check(err)
		// Skip lines with negative complexity
		if com < 0 { continue }
		co := new(Cmplx)
		co.Pos = pos
		co.Cm = com
		co.Chr = chr
		chrCmplx[chr] = append(chrCmplx[chr], *co)
	}
	file.Close()
	return chrCmplx
}
