package core

import (
	"bytes"
	"errors"
	"log"
	"text/template"
	"time"

	"github.com/jmorgan1321/SpaceRep/internal/utils"
)

type TmplMap map[string]*template.Template
type ScopeTmplMap map[string]TmplMap

// ScopedTmplMap is used to create a scoped map of templates like so:
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
var ScopedTmplMap map[string]TmplMap

func init() {
	ScopedTmplMap = map[string]TmplMap{}
}

// Displays allow SpaceRep to the communicate with user, and encapsulate
// the presentation of information to the user.
//
// When a new display is created it should be added into factory.DFE
//
type Display interface {
	// Name should always match the filename of the card.
	Name() string

	// Allows displays to be assoicated with templates.
	SetTmpl(int)
	// Tmpl should return the name of the template used to render a
	// display.
	Tmpl() string

	// Displays know how to create corresponding info data.
	CreateInfo(name string) []Info
	SetInfo(*Info)

	Type() string
}

type Info struct {
	File string
	Set  string `json:"-"`
	// Type is interpreted by Displays to mean different things
	Type int
	// Count stores how many times the user got the correct result on this card.
	//
	// Depending on the bucket the user has to reach a count of 'N' in-order to
	// move this card into the next highest bucket.  A count below zero causes
	// this card to be moved into a lower bucket.
	//
	Count               int
	Bucket              Bucket
	FirstSeen, LastSeen time.Time
}

type Card struct {
	Display
	Info
}

func (c *Card) IncCount() {
	c.Count = utils.Clamp(c.Count+1, 0, c.Bucket.GetMaxCount())
}

func (c *Card) DecCount() {
	c.Count--
}

func (c *Card) Tmpl() string {
	return c.Display.Tmpl()
}

func (c *Card) UpdateBucket() {
	if c.Count >= c.Bucket.GetMaxCount() {
		c.Count = 0
		c.Bucket = c.Bucket.NextBucket()
	}
	if c.Count < 0 {
		c.Count = 0
		c.Bucket = c.Bucket.PrevBucket()
	}
}

// Render takes the display and presents it as an html string.
// An error is return if the card fails to render (ie, bad template)
// or if it can't find a template to render (ie, bad template name).
func (c *Card) Render() (string, error) {
	scopes := []string{c.Set, "html"}

	for _, scope := range scopes {
		if tmap, found := ScopedTmplMap[scope]; found {
			if tmpl, found := tmap[c.Tmpl()]; found {
				var html bytes.Buffer
				if err := tmpl.Execute(&html, c.Display); err != nil {
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
