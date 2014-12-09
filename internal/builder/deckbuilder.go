package builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jmorgan1321/SpaceRep/displays/factory"
	"github.com/jmorgan1321/SpaceRep/internal/core"
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (b *Builder) getDecks() ([]string, error) {
	decks := []string{}
	if b.deck != "all" {
		decks = append(decks, strings.ToLower(b.deck))
		// TODO: don't need to this here, if we check when user sets flag.
		if found, _ := exists(b.rootPath() + b.deck); !found {
			return nil, errors.New("deck does not exist: " + b.deck)
		}
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
	return decks, nil
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
//  	"ppc":
//			{}
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

type deckInfo struct {
	Display string
	Info    []*core.Info
}

func getDeckInfo(path string) (*deckInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	var di deckInfo
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&di)
	if err != nil {
		return nil, err
	}

	return &di, nil
}

func (b *Builder) rootPath() string {
	path := ""
	if b.dir != "" {
		path += (b.dir + "/")
	}
	return path + "html/decks/"
}

func getDataFromDisk(path, display string) ([]core.Display, error) {
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

		d, err := factory.DFE(display)
		if err != nil {
			panic(err)
		}
		decoder.Decode(d)
		alldata = append(alldata, d)

		return nil
	})

	return alldata, nil
}

func updateBuckets(cards []*core.Card) []*core.Card {
	out := []*core.Card{}
	for _, c := range cards {
		c.UpdateBucket()
		out = append(out, c)
	}
	return out
}

func makeCards(set string, info []*core.Info, data []core.Display) []*core.Card {
	// throw the info's in a map
	fileMap := map[string]bool{}
	for _, i := range info {
		fileMap[i.File] = false
	}

	cards := []*core.Card{}
	displayMap := map[string]core.Display{}
	// Figure out which displays are new, by checking against the fileMap.
	// A display is new if it isn't in the fileMap.
	for _, d := range data {
		displayMap[d.Name()] = d
		if _, found := fileMap[d.Name()]; found {
			fileMap[d.Name()] = true
		} else {
			for _, i := range d.CreateInfo(d.Name()) {
				i.Set = set
				d.SetTmpl(int(i.Type))
				c := &core.Card{Info: i, Display: d}
				c.Display.SetInfo(&c.Info)
				cards = append(cards, c)
			}
		}
	}

	// Only add cards that haven't been deleted,
	//	by checking against fileMap.
	for _, i := range info {
		if inBoth := fileMap[i.File]; inBoth {
			i.Set = set
			displayMap[i.File].SetTmpl(int(i.Type))
			c := &core.Card{Info: *i, Display: displayMap[i.File]}
			c.Display.SetInfo(&c.Info)
			cards = append(cards, c)
		}
	}

	return cards
}

// TODO: clean this up or move it.  Exported for savehandler
var SaveDeck = writeDeckInfo

func writeDeckInfo(path, display string, cards []*core.Card) {
	info := []*core.Info{}
	for _, c := range cards {
		info = append(info, &c.Info)
	}

	d := deckInfo{Display: display, Info: info}
	b, _ := json.MarshalIndent(d, "", "\t")

	f, err := os.Open(path)
	if err != nil {
		panic("file read error with: " + path)
	}
	defer f.Close()

	f.Write(b)
}

func (b *Builder) LoadDeck() (*core.Deck, error) {
	sets, err := b.getDecks()
	if err != nil {
		return nil, err
	}
	b.loadTemplates(sets)

	deck := &core.Deck{}
	for _, set := range sets {
		root := b.rootPath() + set
		info, err := getDeckInfo(root + "/cards/cards.info")
		if err != nil {
			return nil, err
		}

		data, err := getDataFromDisk(root+"/cards", info.Display)
		if err != nil {
			return nil, err
		}

		cards := makeCards(set, info.Info, data)
		updateBuckets(cards)

		writeDeckInfo(root+"/cards/cards.info", info.Display, cards)

		for _, c := range cards {
			deck[c.Bucket] = append(deck[c.Bucket], c)
		}
	}

	return deck, nil
}

type Card interface {
	// Info
	Display

	// UpdateBucket()
}

// type Info interface {
// 	File() string
// 	Set() string
// 	// Type() int
// 	// Count() int
// 	// Bucket() int
// }

type Display interface {
	SetInfo(Info)
	Info() Info
	Type() string // type of display.  ie, "basic"
	Name() string // The thing we're displaying
	Tmpl() string

	Clone(Info) Display
	Render(core.ScopeTmplMap) (string, error)
	// CreateInfo(name string) []Info
}

type Info struct {
	File, Set     string
	Type          int
	Count, Bucket int
}

type BasicDisplay struct {
	Stat                          Info
	Word, Image, Desc, Hint, Comp string
}

func (b *BasicDisplay) Name() string   { return b.Word }
func (b *BasicDisplay) Type() string   { return "basic" }
func (b *BasicDisplay) Tmpl() string   { return "thisdoesx" }
func (b *BasicDisplay) SetInfo(i Info) { b.Stat = i }
func (b *BasicDisplay) Info() Info     { return b.Stat }
func (b *BasicDisplay) Clone(i Info) Display {
	return &BasicDisplay{
		Stat:  i,
		Word:  b.Word,
		Image: b.Image,
		Desc:  b.Desc,
		Hint:  b.Hint,
		Comp:  b.Comp,
	}
}
func (b *BasicDisplay) Render(s core.ScopeTmplMap) (string, error) {
	scopes := []string{b.Info().Set, "html"}

	for _, scope := range scopes {
		if tmap, found := s[scope]; found {
			if tmpl, found := tmap[b.Tmpl()]; found {
				var html bytes.Buffer
				if err := tmpl.Execute(&html, b); err != nil {
					return "", err
				}
				return html.String(), nil
			}
		} else {
			log.Fatal("deck not found in scope: ", scope)
		}
	}

	return "", errors.New("couldn't find template to render")

}

type DispHolder struct {
	File string
	Disp Display
}

func makeCards2(set string, info []Info, hldr []DispHolder) []Card {
	// throw the info's in a map
	fileMap := map[string]bool{}
	for _, i := range info {
		fileMap[i.File] = false
	}

	cards := []Card{}
	displayMap := map[string]Display{}
	// Figure out which displays are new, by checking against the fileMap.
	// A display is new if it isn't in the fileMap.
	for _, h := range hldr {
		d := h.Disp
		displayMap[h.File] = d
		if _, found := fileMap[h.File]; found {
			fileMap[h.File] = true
		} else {
			nc, _ := TempDFE(d, h.File)
			cards = append(cards, nc...)
		}
	}

	// Only add cards that haven't been deleted,
	//	by checking against fileMap.
	for _, i := range info {
		if inBoth := fileMap[i.File]; inBoth {
			d := displayMap[i.File].Clone(i)
			cards = append(cards, d)
		}
	}

	return cards
}

func CreateCards(d Display, file string) []Card {
	return []Card{
		d.Clone(Info{File: file, Type: 0}),
		d.Clone(Info{File: file, Type: 1}),
	}
}

// TODO: change the case labels to be names of displays that
//       get stored in a cards.info file.  ie, "basic"
func TempDFE(d Display, file string) ([]Card, error) {
	switch d.Type() {
	case "basic":
		return CreateCards(d, file), nil
	}
	return nil, errors.New("unknown display type passed in: " + d.Type())
}
