package mau

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type GO    map[string]bool
type SymGO map[string]GO
type GidGO map[int]GO

// getGeneSym reads a geneInfo file and returns a map of gene ids and gene symbols.
func getGeneSym(file *os.File) map[int]string {
	var str1, str2 []string
	const numCol = 16
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	Check(err)
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
		Check(err)
		geneSym[i] = str2[2]
	}
	return geneSym
}

var gidGO GidGO // Maps gene ids to sets of GO-accessions.

// getGidGO reads a gene2go file and returns a map of gene ids and sets of GO-accessions.
// It also constructs a map of GO-accessions and their descriptions.
func getGidGO(file *os.File) GidGO {
	var str1, str2 []string
	const numCol = 8

	if gidGO != nil { return gidGO }

	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory
	gidGO = make(GidGO)
	goDescr = make(map[string]string)
	// Iterate over lines
	for _, s := range str1 {
		// Skip empty lines and commend lines
		if len(s) == 0 || s[0:1] == "#" { continue }
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) < numCol {
			fmt.Fprintf(os.Stderr, "Warning[GetGidGO]: Expect %d columns, but find only %d in %s\n", numCol, len(str2), s)
			continue
		}
		i, err := strconv.Atoi(str2[1])
		if err != nil {
			log.Fatal(err)
		}
		goAcc := gidGO[i]
		if goAcc == nil {
			goAcc = make(map[string]bool)
			gidGO[i] = goAcc
		}
		goAcc[str2[2]] = true
		goDescr[str2[2]] = str2[5]
	}
	return gidGO
}

var goDescr map[string]string

// GOdescr takes as input a gene-info file and a gene2go file and returns a map of GO-accessions and their descriptions.
func GOdescr() map[string]string {
	if goDescr == nil {
		log.Fatal("GO-descriptions not available; need to call getGidGO first")
	}

	return goDescr
}

// GetSymGO takes as input a gene-info file and a gene2go file and returns a map of gene symbols and sets of GO-accessions.
func GetSymGO(infoFile, ggFile string) SymGO {

	var infoF, geneF *os.File
	var err error
	
	infoF, err = os.Open(infoFile)
	Check(err)
	geneF, err = os.Open(ggFile)
	Check(err)

	geneSym := getGeneSym(infoF)
	gidGO  := getGidGO(geneF)
	symGO := make(SymGO)
	for ge, sy := range geneSym {
		ga := gidGO[ge]
		goAcc := symGO[sy]
		if goAcc == nil {
			goAcc = make(GO)
			symGO[sy] = goAcc
		}
		for k, _ := range ga {
			goAcc[k] = true
		}
	}
	return symGO
}

