package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var directory string
var noColor bool

func main() {
	flag.StringVar(&directory, "case-dir", "", "Directory for the test cases")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")

	flag.Parse()

	colorString := strconv.FormatBool(noColor)

	cmd := exec.Command("go", "test", fmt.Sprintf("-no-color=%s", colorString), fmt.Sprintf("-case-dir=%s", directory))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
