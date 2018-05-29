package mau

import (
	"strconv"
	"strings"
	"io/ioutil"
	"os"
	"bufio"
)
// 0    1               2       3       4       5       6       7
// 585	NR_046018	chr1	+	11873	14409	14409	14409	3	11873,12612,13220,	12227,12721,14409,	0	DDX11L1	unk	unk	-1,-1,-1,
type Gene struct {
	Start int         // column 4
	End   int         // column 5
	Sym   string      // column 12
}

type ChrGene map[string][]Gene

// GetChrRef reads a refGene file and returns a map of chromosomes to genes.
func GetChrGene(fileName string) ChrGene {
	var str1, str2 []string
	var chr string
	var i, j int
	var file *os.File
	var err error

	file, err = os.Open(fileName)
	Check(err)
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	Check(err)
	// Split data into lines
	str1 = strings.Split(string(data), "\n")
	// Allocate memory 
	chrGene := make(ChrGene)
	// Iterate over lines
	for _, s := range str1 {
		// Split line into columns
		str2 = strings.Split(s, "\t")
		if len(str2) != 16 {
			continue
		}
		chr = str2[2]
		i, err = strconv.Atoi(str2[4])
		Check(err)
		j, err = strconv.Atoi(str2[5])
		Check(err)
		gene := new(Gene)
		gene.Start = i
		gene.End   = j
		gene.Sym   = str2[12]
		chrGene[chr] = append(chrGene[chr], *gene)
	}
	file.Close()
	return chrGene
}

