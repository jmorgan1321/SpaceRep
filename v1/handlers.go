package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/jmorgan1321/SpaceRep/sets/work"
// 	"net/http"
// )

// type Response map[string]interface{}

// func (r Response) String() (s string) {
// 	b, err := json.Marshal(r)
// 	if err != nil {
// 		s = ""
// 		return
// 	}
// 	s = string(b)
// 	return
// }

// // TODO: just remove this for the card itself
// type tmplWrapper struct {
// 	Set string
// 	*work.Display
// }

// func cardIndexHandler(rw http.ResponseWriter, r *http.Request) {
// }

// func reviewHandler(rw http.ResponseWriter, r *http.Request) {
// 	// create new flashcard html
// 	var html bytes.Buffer
// 	key := g_currentCard.Type.String()
// 	// key := g_currentCard.Info.Tmpl
// 	g_env.TmplMap[key].Execute(&html, g_currentCard)

// 	// decode web client's message
// 	decoder := json.NewDecoder(r.Body)
// 	var t struct{ Status string }
// 	decoder.Decode(&t)
// 	fmt.Printf("1: %#v\n", t)

// 	// update card count based on user response
// 	fmt.Println(t.Status)
// 	switch t.Status {
// 	case "Accept":
// 		g_currentCard.IncCount()
// 	case "Forgot":
// 		g_currentCard.DecCount()
// 	}

// 	// send next card's html back to webpage
// 	rw.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(rw, Response{
// 		"success": true,
// 		"message": "Hello!",
// 		"newCard": html.String(),
// 	})
// }

// func saveHandler(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(rw, Response{
// 		"success": true,
// 		"message": "Hello!",
// 	})

// 	// saveDeck(g_env.Cards)
// }

// func submitHandler(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(rw, Response{"success": true, "message": "Hello!"})

// 	// decoder := json.NewDecoder(r.Body)

// 	// var t work.Display
// 	// decoder.Decode(&t)

// 	// // write out new card data
// 	// f, err := os.Create("html/" + g_currentCard.Info.Set + "/cards/" + t.Word + ".data")
// 	// if err != nil {
// 	//  panic(err)
// 	// }
// 	// defer f.Close()

// 	// b, err := json.MarshalIndent(t, "", "\t")
// 	// if err != nil {
// 	//  panic(err)
// 	// }
// 	// f.Write(b)
// }
