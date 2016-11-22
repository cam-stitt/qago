package main

import (
	"flag"
	"os"

	"github.com/cam-stitt/qago"
)

var directory string
var noColor bool

func main() {
	flag.StringVar(&directory, "case-dir", "./fixtures", "Directory for the test cases")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")

	flag.Parse()

	suite := qago.Suite{
		Directory: directory,
		NoColor:   noColor,
	}

	os.Exit(suite.Run())
}
