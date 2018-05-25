package macle2go

import (
	"strconv"
	"strings"
	"log"
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
)

type Macle struct {
	Chr string
	Pos int
	Cm  float64
}

func colCount(col []string) int {
	c := 0
	for _, s := range col {
		if len(s) > 0 { c++ }
	}
	return c
}

// compWinLen: Compute length of sliding window
func compWinLen(data []string) int {
	var str []string
	
	str = strings.Split(data[0], "\t")
	s, err := strconv.ParseFloat(str[1], 64)
	if err != nil { log.Fatal(err) }
	w := int(2 * s)

	return w
}

// GetMacle: Process string data obtained from a macle file and return a slice of macle entries
func GetMacle(file *os.File) []Macle {
	var arr []Macle
	var str1, str2 []string
	var f float64
	const numCol = 3

	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory 
	arr = make([]Macle, len(str1))
	n := 0
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
		f, err = strconv.ParseFloat(str2[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		// Skip lines with negative complexity
		if f < 0 {
			continue
		}
		arr[n].Cm = f
		arr[n].Chr = str2[0]
		arr[n].Pos, err = strconv.Atoi(str2[1])
		if err != nil {
			log.Fatal(err)
		}
		n++
	}

	return arr[:n]
}
