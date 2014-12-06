package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	flag.Usage = usage
}

func main() {
	fmt.Println("hello")
	defer fmt.Println("goodbye")

	flag.Parse()
}

func usage() {
	msg := `
    Desc:
    %s does things.

    Usage:    %s [options]

    The options are listed below in the following format:

        -option=default value:  description

    Options:
`
	prog := filepath.Base(os.Args[0])
	fmt.Printf(msg+"\n", prog, prog)
	flag.PrintDefaults()
}
