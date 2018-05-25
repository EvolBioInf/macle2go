package macle2go

import (
	"strconv"
	"strings"
	"log"
	"io/ioutil"
	"os"
	"bufio"
)
// 0    1               2       3       4       5       6       7
// 585	NR_046018	chr1	+	11873	14409	14409	14409	3	11873,12612,13220,	12227,12721,14409,	0	DDX11L1	unk	unk	-1,-1,-1,
type RefGene struct {
	Chr   string      // column 2
	Start int         // column 4
	End   int         // column 5
}

// GetMacle: Process string data obtained from a gene2go file and return a slice of gene2go entries
func GetRefGene(file *os.File) []RefGene {
	var arr []RefGene
	var str1, str2 []string
	var i, j int
	var err error

	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory 
	arr = make([]RefGene, len(str1))
	n := 0
	// Iterate over lines
	for _, s := range str1 {
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) != 16 {
			continue
		}
		arr[n].Chr = str2[2]
		i, err = strconv.Atoi(str2[4])
		if err != nil {
			log.Fatal(err)
		}
		j, err = strconv.Atoi(str2[5])
		arr[n].Start = i
		arr[n].End   = j
		n++
	}
	return arr[:n]
}
