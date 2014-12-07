package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmorgan1321/SpaceRep/v1/displays/basic"
	"github.com/jmorgan1321/SpaceRep/v1/internal/builder"
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
)

func init() {
	flag.Usage = usage
}

func main() {
	doneCh := make(chan struct{})

	fmt.Println("hello")
	defer fmt.Println("goodbye")

	flag.Parse()

	_, err := builder.New(
		builder.DFE(dfe),
		builder.Deck("ppc"),
	).LoadDeck()

	if err != nil {
		log.Fatal("failed to load deck: ", err)
	}

	// spin up server thread
	go func() {
		router := NewRouter()
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./html/static")))
		http.ListenAndServe(":8080", router)
		// log.Fatal(http.ListenAndServe(":8080", router))
	}()

	// // spin up presentation thread
	// go func() {
	// 	for _, b := range deck {
	// 		for _, c := range b {
	// 			fmt.Printf("\t%#v\n", c)
	// 		}
	// 	}
	// 	doneCh <- struct{}{}
	// }()

	<-doneCh
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

func dfe(s string) core.Display {
	switch s {
	case "ppc", "git":
		return &basic.Display{}
	}
	return nil
}
