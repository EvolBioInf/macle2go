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
// 0            1       2       3               4               5       6       7
// tax_id	GeneID	GO_ID	Evidence	Qualifier	GO_term	PubMed	Category
type Gene2go struct {
	GeneId   int     // column 1
	GoId     string  // column 2
	GoTerm   string  // column 5
	Category string  // column 7
}

// GetMacle: Process string data obtained from a gene2go file and return a slice of gene2go entries
func GetGene2go(file *os.File) []Gene2go {
	var arr []Gene2go
	var str1, str2 []string
	const numCol = 8

	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory 
	arr = make([]Gene2go, len(str1))
	n := 0
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines and commend lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) < numCol {
			fmt.Fprintf(os.Stderr, "Warning[GetGene2Go]: Expect %d columns, but find only %d in %s\n", numCol, len(str2), s)
			continue
		}
		i, err := strconv.Atoi(str2[1])
		if err != nil {
			log.Fatal(err)
		}
		arr[n].GeneId   = i
		arr[n].GoId     = str2[2]
		arr[n].GoTerm   = str2[5]
		arr[n].Category = str2[7]
		n++
	}
	return arr[:n]
}
