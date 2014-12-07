package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
)

type Builder struct {
	dir  string      // if dir != "" then 'html' dir is located (from root) in dir
	deck string      // if != "all" then load try to load that specific deck
	dfe  FactoryFunc // DFEExtention
}

type Option func(*Builder)

func New(opts ...Option) *Builder {
	bldr := &Builder{
		deck: "all",
	}

	for _, opt := range opts {
		opt(bldr)
	}

	return bldr
}

func Deck(s string) Option {
	return func(b *Builder) {
		b.deck = s
	}
}

// A FactoryFunc is used to extend deck builder with custom
// displays.
type FactoryFunc func(string) core.Display

// DFE set the builder's DisplayFactoryExtention.
func DFE(f FactoryFunc) Option {
	return func(b *Builder) {
		b.dfe = f
	}
}

func (b *Builder) getDecks() []string {
	decks := []string{}
	if b.deck != "all" {
		decks = append(decks, strings.ToLower(b.deck))
	} else {
		filepath.Walk(b.rootPath(), func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() || fi.Name() != "cards" {
				return nil
			}

			deckName := filepath.Base(filepath.Dir(path))

			if deckName == "html" {
				log.Println("ignoring invalid deck name: html")
				return nil
			}

			decks = append(decks, strings.ToLower(deckName))
			return nil
		})
	}
	return decks
}

func (b *Builder) getDisplay(s string) (core.Display, error) {
	var d core.Display
	if b.dfe != nil {
		d = b.dfe(strings.ToLower(s))
	}
	if d == nil {
		return nil, errors.New("unknown display type passed in: " + s)
	}
	return d, nil
}

func (b *Builder) LoadDeck() (*core.Deck, error) {
	decks := b.getDecks()

	b.loadTemplates(decks)

	alldata := []core.Display{}
	for _, deck := range decks {
		alldata = append(alldata, b.readCardDataFromDisk(deck)...)
	}

	if len(alldata) == 0 {
		return nil, errors.New("No data read from disk.")
	}

	// Associate display with info (d.Name() must match info.File)
	dataMap := map[string]core.Display{}
	for _, d := range alldata {
		dataMap[d.Name()] = d
	}

	// read in .card files (to hash, compare against), only create new cards
	allinfo := []*core.Info{}
	for _, deck := range decks {
		allinfo = append(allinfo, b.readCardInfoFromDisk(deck)...)
	}

	// create and return cards and card buckets
	deck := &core.Deck{}
	for _, c := range allinfo {
		if c.Count >= c.Bucket.GetMaxCount() {
			c.Count = 0
			c.Bucket = c.Bucket.NextBucket()
		}
		if c.Count < 0 {
			c.Count = 0
			c.Bucket = c.Bucket.PrevBucket()
		}

		// set display type here
		d := dataMap[c.File]
		d.SetType(int(c.Type))
		deck[c.Bucket] = append(deck[c.Bucket], &core.Card{Info: *c, Display: d})
	}

	return deck, nil
}

var tmplIndex = `
<div class="front">
    {{template "front" .}}
</div>
<div class="back">
    {{template "back" .}}
</div>
`

// loadTemplates Creates a scoped map of templates like so:
// 	{
//		"html":
// 			{ "thisdoesx": tmpl1, "xdoesthis": tmpl2 }
//  	"ppc": {}
//  	"git":
//			{ "thisdoesx": tmpl3 }
// }
//
// This allows local templates to override default ones
// and for new display types to use their own templates.
//
func (b *Builder) loadTemplates(decks []string) error {
	tmpl, _ := template.New("base").Parse(tmplIndex)

	getTemplateFunc := func(deck string) func(path string, fi os.FileInfo, err error) error {
		return func(path string, fi os.FileInfo, err error) error {
			shouldScan := strings.HasSuffix(path, ".tmpl") && !fi.IsDir()
			if !shouldScan {
				return nil
			}

			// read contents of template file into string
			b, _ := ioutil.ReadFile(path)
			tmplText := string(b)

			// store template in TmplMap.
			tmplName := fi.Name()[:len(fi.Name())-len(".tmpl")]
			tmpl2, err := tmpl.Clone()
			if err != nil {
				fmt.Println("tmpl:", path)
				log.Fatal("failed to clone template:", err)
			}
			tmpl2, err = tmpl2.Parse(tmplText)
			if err != nil {
				fmt.Println("tmpl:", path)
				log.Fatal("failed to parse template:", err)
			}
			core.ScopedTmplMap[deck][tmplName] = tmpl2

			return nil
		}
	}

	// iterate over all decks and find their templates
	for _, deck := range decks {
		tmplDir := b.rootPath() + deck + "/tmpl"
		core.ScopedTmplMap[deck] = core.TmplMap{}

		filepath.Walk(tmplDir, getTemplateFunc(deck))
	}

	// manually add the 'html/tmpl' dir
	tmplDir := strings.Replace(b.rootPath(), "decks", "tmpl", 1)
	core.ScopedTmplMap["html"] = core.TmplMap{}
	filepath.Walk(tmplDir, getTemplateFunc("html"))

	return nil
}

func (b *Builder) readCardDataFromDisk(deck string) []core.Display {
	alldata := []core.Display{}
	path := b.rootPath() + deck + "/cards"

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

		deck := filepath.Base(filepath.Dir(filepath.Dir(path)))
		d, err := b.getDisplay(deck)
		if err != nil {
			panic(err)
		}
		decoder.Decode(d)
		alldata = append(alldata, d)

		return nil
	})

	return alldata
}

func (b *Builder) readCardInfoFromDisk(deck string) []*core.Info {
	allinfo := []*core.Info{}
	path := b.rootPath() + deck + "/cards/cards.info"

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

	for _, info := range allinfo {
		info.Set = deck
	}

	return allinfo
}

func (b *Builder) rootPath() string {
	path := ""
	if b.dir != "" {
		path += (b.dir + "/")
	}
	return path + "html/decks/"
}
