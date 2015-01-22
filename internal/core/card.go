package core

import (
	"bytes"
	"errors"
	"log"
	"text/template"
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

// type Card interface {
// 	// Info
// 	Tmpl() string
// 	Bucket() Bucket
// 	UpdateBucket()
// 	Set() string
// 	Stats() *Info
// }

type Card struct {
	Info
	Data map[string]string
}

type Info struct {
	File string
	Set  string `json:"Set"`
	// Type   int
	Count  int
	Bucket Bucket `json:"Bucket"`
	Tmpl   string `json:"Tmpl"`
}

func (i *Info) UpdateBucket() {
	if i.Count >= i.Bucket.GetMaxCount() {
		i.Count = 0
		i.Bucket = i.Bucket.NextBucket()
	}
	if i.Count < 0 {
		i.Count = 0
		i.Bucket = i.Bucket.PrevBucket()
	}
}
func (i *Info) IncCount() {
	i.Count++
}
func (i *Info) DecCount() {
	i.Count--
}

func Render(s ScopeTmplMap, c *Card) (string, error) {
	scopes := []string{c.Set, "html"}

	for _, scope := range scopes {
		if tmap, found := s[scope]; found {
			if tmpl, found := tmap[c.Tmpl]; found {
				var html bytes.Buffer
				if err := tmpl.Execute(&html, c); err != nil {
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
