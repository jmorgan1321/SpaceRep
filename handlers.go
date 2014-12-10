package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmorgan1321/SpaceRep/internal/core"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func cardIndexHandler(rw http.ResponseWriter, r *http.Request) {
}

func reviewHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["value"]

	// update card count based on user response
	switch value {
	case "accept":
		g_currCard.Stats().IncCount()
	case "forgot":
		g_currCard.Stats().DecCount()
	}

	// // Update card seen information
	// if g_currCard.FirstSeen.IsZero() {
	// 	g_currCard.FirstSeen = time.Now()
	// }
	// g_currCard.LastSeen = time.Now()

	// create new flashcard html
	html, err := core.Render(core.ScopedTmplMap, g_currCard)
	if err != nil {
		log.Fatal(err)
	}

	// send next card's html back
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"success": true,
		"message": "Hello!",
		"newCard": html,
	})

	// signal that we're ready for the next card
	nextCardCh <- struct{}{}
}

func saveHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"success": true,
		"message": "Hello!",
	})

	saveDeck(g_env.Deck)
}
