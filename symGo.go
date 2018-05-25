package macle2go

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)



// getGeneSym reads a geneInfo file and returns a map of gene ids and gene symbols.
func getGeneSym(file *os.File) map[int]string {
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
	geneSym := make(map[int]string)
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines & comment lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) < numCol {
			fmt.Fprintf(os.Stderr, "Warning[GetGeneSym]: Expect %d columns, but find only %d in %s\n", numCol, len(str2), s)
			continue
		}
		i, err := strconv.Atoi(str2[1])
		if err != nil {
			log.Fatal(err)
		}
		geneSym[i] = str2[2]
	}
	return geneSym
}

var geneGo map[int]map[string]bool // Maps gene ids to sets of GO-accessions.

// getGeneGo reads a gene2go file and returns a map of gene ids and sets of GO-accessions.
// It also constructs a map of GO-accessions and their descriptions.
func getGeneGo(file *os.File) map[int]map[string]bool {
	var str1, str2 []string
	const numCol = 8

	if geneGo != nil { return geneGo }

	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory
	geneGo = make(map[int]map[string]bool)
	goDescr = make(map[string]string)
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
		goAcc := geneGo[i]
		if goAcc == nil {
			goAcc = make(map[string]bool)
			geneGo[i] = goAcc
		}
		goAcc[str2[2]] = true
		goDescr[str2[2]] = str2[5]
	}
	return geneGo
}

var goDescr map[string]string

// GetGoDescr takes as input a gene-info file and a gene2go file and returns a map of GO-accessions and their descriptions.
func GetGoDescr(infoF *os.File, geneF *os.File) map[string]string {
	if goDescr == nil {
		getGeneGo(geneF)
	}
	return goDescr
}

// GetSymGo takes as input a gene-info file and a gene2go file and returns a map of gene symbols and sets of GO-accessions.
func GetSymGo(infoF *os.File, geneF *os.File) map[string]map[string]bool {

	geneSym := getGeneSym(infoF)
	geneGo  := getGeneGo(geneF)

	symGo := make(map[string]map[string]bool)
	for ge, sy := range geneSym {
		ga := geneGo[ge]
		goAcc := symGo[sy]
		if goAcc == nil {
			goAcc = make(map[string]bool)
			symGo[sy] = goAcc
		}
		for k, _ := range ga {
			goAcc[k] = true
		}
	}
	return symGo
}
