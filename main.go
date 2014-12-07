package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmorgan1321/SpaceRep/v1/internal/builder"
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
	"github.com/jmorgan1321/SpaceRep/v1/internal/env"
)

var (
	set     = flag.String("deck", "all", "which deck(s) to load cards from")
	browser = flag.Bool("browser", false, "attempt to open a browser")
	port    = flag.String("port", ":8080", "which port to run card server on")
)

var (
	g_env      *env.Env
	g_currCard *core.Card
	nextCardCh chan struct{}
)

func init() {
	flag.Usage = usage
}

func main() {
	doneCh := make(chan struct{})
	nextCardCh = make(chan struct{})

	flag.Parse()

	deck, err := builder.New(
		builder.DFE(dfe),
		builder.Deck(*set),
	).LoadDeck()

	if err != nil {
		log.Fatal("failed to load deck: ", err)
	}

	g_env = env.New(
		env.Deck(deck),
		env.DistributionFunc(core.StandardDistribution),
	)

	// spin up server thread
	go func() {
		router := NewRouter()
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./html")))
		log.Fatal(http.ListenAndServe(*port, router))
	}()

	// spin up presentation thread
	go func() {
		for card := range g_env.Deck.GetCards(g_env.Seed, g_env.Distributions) {
			g_currCard = card
			<-nextCardCh
		}
	}()

	if *browser {
		openBrowser("http://localhost" + *port)
	} else {
		fmt.Println("open browser to http://localhost" + *port)
	}

	<-doneCh
}

func usage() {
	msg := `
    Desc:
    %s is a spaced repetition system (srs) designed to optimize
    learning and retention by presenting flashcards at specific
    intervals!

    The frequency a flashcard is shown depends on how
    many times the cards was remembered or forgotten and how long
    it's been since the card was last seen.


    Usage:    %s [options]

    The options are listed below in the following format:

        -option=default value:  description

    Options:
`
	prog := filepath.Base(os.Args[0])
	fmt.Printf(msg+"\n", prog, prog)
	flag.PrintDefaults()
}

// TODO: move or make a member of Deck?
func saveDeck(deck *core.Deck) {
	sets := map[string][]*core.Card{}
	for _, bucket := range deck {
		for _, card := range bucket {
			sets[card.Set] = append(sets[card.Set], card)
		}
	}

	for set, cards := range sets {
		fmt.Println("saving deck:", set)
		saveCards(cards)
	}
}

func saveCards(cards []*core.Card) {
	allinfo := []*core.Info{}

	for _, card := range cards {
		allinfo = append(allinfo, &card.Info)
	}

	// rewrite .cards file
	b, err := json.MarshalIndent(allinfo, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	path := "html/decks/" + cards[0].Set + "/cards/cards.info"
	f, err := os.Create(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	f.Write(b)
}
