package basic

import "github.com/jmorgan1321/SpaceRep/internal/core"

type Card struct {
	core.Info
	Word, Image string
	Desc, Hint  string
	Comp        string
}

// Card interface
func (c *Card) Name() string { return c.Word }
func (c *Card) Type() string { return "basic" }
func (c *Card) Clone(i core.Info) core.Card {
	return &Card{
		Info: core.Info{
			File:  i.File,
			S:     i.S,
			Type:  i.Type,
			Count: i.Count,
			B:     i.B,
		},
		Word:  c.Word,
		Image: c.Image,
		Desc:  c.Desc,
		Hint:  c.Hint,
		Comp:  c.Comp,
	}
}

func (c *Card) Tmpl() string {
	return Type(c.Info.Type).String()
}

func CreateCardsFromTemplate(c core.Card, i core.Info) []core.Card {
	return []core.Card{
		c.Clone(core.Info{File: i.File, S: i.S, Type: int(DescCard)}),
		c.Clone(core.Info{File: i.File, S: i.S, Type: int(WordCard)}),
	}
}
