package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/evolbioinf/macle2go"
)

func main() {
	var gc = flag.Float64("g", 0.5, "GC-content")
	var len = flag.Float64("l", 1000000, "Genome-length")
	var win = flag.Float64("w", 5000, "Window-length")
	var prob = flag.Float64("p", 0.05, "Probability")
	var ver = flag.Bool("v", false, "Version")

	flag.Parse()
	if *ver == true {
		fmt.Println("Version 0.1")
		os.Exit(0)
	}
	fmt.Printf("%v\t%v\n", *prob, macle2go.Quantile(*gc, *len, *win, *prob))
}
