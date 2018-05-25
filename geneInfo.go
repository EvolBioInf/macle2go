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
// #tax_id	GeneID	Symbol	LocusTag	Synonyms	dbXrefs	chromosome	map_location	description	type_of_gene	Symbol_from_nomenclature_authority	Full_name_from_nomenclature_authority	Nomenclature_status	Other_designations	Modification_date	Feature_type
type GeneInfo struct {
	GeneId   int     // column 1
	Sym      string  // column 2
}

// GetMacle: Process string data obtained from a gene2go file and return a slice of gene2go entries
func GetGeneInfo(file *os.File) []GeneInfo {
	var arr []GeneInfo
	var str1, str2 []string
	const numCol = 16
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory 
	arr = make([]GeneInfo, len(str1))
	n := 0
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines & comment lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) < numCol {
			fmt.Fprintf(os.Stderr, "Warning[GetGeneInfo]: Expect %d columns, but find only %d in %s\n", numCol, len(str2), s)
			continue
		}
		i, err := strconv.Atoi(str2[1])
		if err != nil {
			log.Fatal(err)
		}
		arr[n].GeneId   = i
		arr[n].Sym      = str2[2]
		n++
	}
	return arr[:n]
}
