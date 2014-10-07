package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmorgan1321/srs/core"
	"github.com/jmorgan1321/srs/env"
	"github.com/jmorgan1321/srs/sets/work"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var g_env *env.Env
var g_currentCard *core.Card

func init() {
	env.LoadCardsFunc = loadCards
	env.DistributionFunc = core.StandardDistribution

	g_env = env.New()

	cardHtmlTmpl1 := `
<div class="front">
    <p>What does the {{.Comp}} <b>{{.Word}}</b> do?</p>
</div>
<div class="back">
    <p>{{.Desc}}</p>
    <img src="facts/image/{{.Image}}" height="150" width="150" />
    <p>({{.Hint}})</p>
</div>`
	tmpl1, _ := template.New("test").Parse(cardHtmlTmpl1)
	g_env.TmplMap[core.WordCard.String()] = tmpl1

	cardHtmlTmpl2 := `
<div class="front">
    <p>What {{.Comp}} {{.Desc}}?</p>
</div>
<div class="back">
    <p><b>{{.Word}}</b></p>
    <img src="facts/image/{{.Image}}" height="150" width="150" />
    <p>({{.Hint}})</p>
</div>`
	tmpl2, _ := template.New("test").Parse(cardHtmlTmpl2)
	g_env.TmplMap[core.DescCard.String()] = tmpl2
}

type selection struct {
	x    float32
	card *core.Card
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handled")
	w.Write([]byte("Hello root."))

	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

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

func reviewHandler(rw http.ResponseWriter, r *http.Request) {
	// create new flashcard html
	var html bytes.Buffer
	g_env.TmplMap[g_currentCard.Type.String()].Execute(&html, g_currentCard.Display.(*work.Display))

	// decode web client's message
	decoder := json.NewDecoder(r.Body)
	var t struct{ Status string }
	decoder.Decode(&t)
	fmt.Printf("1: %#v\n", t)

	// update card count based on user response
	fmt.Println(t.Status)
	switch t.Status {
	case "Accept":
		g_currentCard.IncCount()
	case "Forgot":
		g_currentCard.DecCount()
	}

	// send next card's html back to webpage
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"success": true,
		"message": "Hello!",
		"newCard": html.String(),
	})
}

func saveHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"success": true,
		"message": "Hello!",
	})

	saveDeck(g_env.Cards)
}

func submitHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{"success": true, "message": "Hello!"})

	decoder := json.NewDecoder(r.Body)

	var t work.Display
	decoder.Decode(&t)

	// write out new card data
	f, err := os.Create("html/facts/cards/" + t.Word + ".data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		panic(err)
	}
	f.Write(b)
}

func openBrowser(url string) {
	cmd := exec.Command("cmd", "/c", "start", url)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", url)
	}
}

func main() {
	doneCh := make(chan bool)
	readyCh := make(chan bool)

	go func() {
		dir := "C:/Users/jmorgan/Sandbox/golang/src/github.com/jmorgan1321/srs/html/"

		http.HandleFunc("/api/submit", submitHandler)
		http.HandleFunc("/api/save", saveHandler)
		// TODO: clean this up.
		http.HandleFunc("/api/review", func(fn http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				// if origin := r.Header.Get("Origin"); origin != "" {
				// 	w.Header().Set("Access-Control-Allow-Origin", origin)
				// }
				// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
				// w.Header().Set("Access-Control-Allow-Credentials", "true")

				fn(w, r)
				readyCh <- true
			}
		}(reviewHandler))
		http.Handle("/", http.FileServer(http.Dir(dir)))
		http.ListenAndServe(":8080", nil)
	}()

	go func() {
		fmt.Printf("[%d %d %d %d]\n",
			len(g_env.Cards[core.Daily]), len(g_env.Cards[core.Weekly]),
			len(g_env.Cards[core.Monthly]), len(g_env.Cards[core.Yearly]))
		fmt.Println(g_env.Distributions)

		ch := getCards(g_env)
		for selection := range ch {
			g_currentCard = selection.card
			presentSelectionToUser(selection)
			<-readyCh
		}
	}()

	createCards()

	<-time.After(1 * time.Second)
	// TODO: make openBrowser flag
	// openBrowser("http://localhost:8080")

	<-doneCh
}

func getCards(env *env.Env) <-chan selection {
	r := rand.New(rand.NewSource(42))
	ch := make(chan selection)

	go func() {
		for {
			// get random distribution selection between [0...1]
			x := r.Float32()

			bucket := core.Daily
			switch {
			case x < env.Distributions[core.Yearly]:
				bucket = core.Yearly
			case x < env.Distributions[core.Monthly]:
				bucket = core.Monthly
			case x < env.Distributions[core.Weekly]:
				bucket = core.Weekly
			}

			// pick random card from bucket
			// TODO: handle empty bucket case:
			if len(env.Cards[bucket]) == 0 {
				continue
			}
			i := r.Intn(len(env.Cards[bucket]))
			ch <- selection{x: x, card: env.Cards[bucket][i]}
		}
	}()

	return ch
}

func presentSelectionToUser(s selection) {
	fmt.Printf("%.2f - %7s - %10s - %2d\n", s.x, s.card.Bucket, s.card.Display.(*work.Display).Word, s.card.Count)
}

// TODO: promote to core?  Do all decks load cards the same way
func ReadCardDataFromDisk(path string) []core.Display {
	alldata := []core.Display{}

	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		shouldScan := strings.HasSuffix(path, ".data") && !fi.IsDir()
		if !shouldScan {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			panic("file read error with: " + path)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)

		var d work.Display
		decoder.Decode(&d)

		alldata = append(alldata, &d)
		return nil
	})

	return alldata
}

func ReadCardInfoFromDisk(path string) []*core.Info {
	allinfo := []*core.Info{}

	f, err := os.Open(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&allinfo)
	if err != nil {
		fmt.Println("err:", err)
	}

	return allinfo
}

func saveDeck(deck [core.BucketCount][]*core.Card) {
	allinfo := []*core.Info{}

	for _, bucket := range deck {
		for _, card := range bucket {
			allinfo = append(allinfo, &card.Info)
		}
	}

	// rewrite .cards file
	b, err := json.MarshalIndent(allinfo, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	path := "html/facts/cards/cards.info"
	f, err := os.Create(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	f.Write(b)
}

func createCards() {
	// read in all .data files
	alldata := ReadCardDataFromDisk("html/facts/cards")
	fmt.Println("\ndata:")
	for _, v := range alldata {
		fmt.Printf("\t%#v\n", v)
	}

	// read in .card files (to hash, compare against), only create new cards
	allinfo := ReadCardInfoFromDisk("html/facts/cards/cards.info")
	fmt.Println("\ninfo:")
	for _, v := range allinfo {
		fmt.Printf("\t%#v\n", v)
	}

	// TODO: make ReadCardInfoFromDisk return [core.numcardtypes]map[string]*core.Info{}
	infoMap := map[string]map[string]*core.Info{
		core.WordCard.String(): map[string]*core.Info{},
		core.DescCard.String(): map[string]*core.Info{},
	}
	for _, info := range allinfo {
		infoMap[info.Type.String()][info.File] = info
	}

	// for each item in alldata, if it doesn't exist in allinfo--create it
	// for all data create .cards
	for _, d := range alldata {
		dd := d.(*work.Display)
		if _, found := infoMap[core.DescCard.String()][dd.Word]; !found {
			allinfo = append(allinfo, &core.Info{File: dd.Word, Type: core.DescCard})
		}
		if _, found := infoMap[core.WordCard.String()][dd.Word]; !found {
			allinfo = append(allinfo, &core.Info{File: dd.Word, Type: core.WordCard})
		}
	}

	// rewrite .cards file
	b, err := json.MarshalIndent(allinfo, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	path := "html/facts/cards/cards.info"
	f, err := os.Create(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	f.Write(b)
}

func loadCards() [core.BucketCount][]*core.Card {
	alldata := ReadCardDataFromDisk("html/facts/cards")
	dataMap := map[string]core.Display{}
	for _, d := range alldata {
		dataMap[d.(*work.Display).Word] = d
	}

	allinfo := ReadCardInfoFromDisk("html/facts/cards/cards.info")

	// create and return cards and card buckets
	deck := [core.BucketCount][]*core.Card{}
	for _, c := range allinfo {
		if c.Count >= c.Bucket.GetMaxCount() {
			c.Count = 0
			c.Bucket = c.Bucket.NextBucket()
		}
		if c.Count < 0 {
			c.Count = 0
			c.Bucket = c.Bucket.PrevBucket()
		}
		deck[c.Bucket] = append(deck[c.Bucket], &core.Card{Info: *c, Display: dataMap[c.File]})
	}

	return deck
}
