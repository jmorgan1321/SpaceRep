package basic

import "github.com/jmorgan1321/SpaceRep/internal/core"

type Card struct {
	Info        *core.Info
	Word, Image string
	Desc, Hint  string
	Comp        string
}

// Card interface
func (c *Card) Name() string { return c.Word }
func (c *Card) Type() string { return "basic" }
func (c *Card) Clone(i core.Info) core.Card {
	return &Card{
		Info: &core.Info{
			File:   i.File,
			Set:    i.Set,
			Type:   i.Type,
			Count:  i.Count,
			Bucket: i.Bucket,
		},
		Word:  c.Word,
		Image: c.Image,
		Desc:  c.Desc,
		Hint:  c.Hint,
		Comp:  c.Comp,
	}
}
func (c *Card) Bucket() core.Bucket {
	return c.Info.GetBucket()
}
func (c *Card) UpdateBucket() {
	c.Info.UpdateBucket()
}
func (c *Card) Set() string {
	return c.Info.Set
}
func (c *Card) Tmpl() string {
	return Type(c.Info.Type).String()
}
func (c *Card) Stats() *core.Info {
	return c.Info
}

func CreateCardsFromTemplate(c core.Card, i core.Info) []core.Card {
	return []core.Card{
		c.Clone(core.Info{File: i.File, Set: i.Set, Type: int(DescCard)}),
		c.Clone(core.Info{File: i.File, Set: i.Set, Type: int(WordCard)}),
	}
}
