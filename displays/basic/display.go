package basic

import "github.com/jmorgan1321/SpaceRep/internal/core"

// Display represent field that we want to display on cards.
type Display struct {
	Word, Image, Desc, Hint string
	Comp                    string
	Typ                     Type
	*core.Info
}

func (d *Display) Name() string {
	return d.Word
}

func (d *Display) SetTmpl(id int) {
	d.Typ = Type(id)
}

func (d *Display) Tmpl() string {
	return d.Typ.String()
}

func (d *Display) Type() string {
	return "basic"
}

func (d *Display) SetInfo(i *core.Info) {
	d.Info = i
}

func (d *Display) CreateInfo(word string) []core.Info {
	return []core.Info{
		core.Info{File: word, Type: int(DescCard)},
		core.Info{File: word, Type: int(WordCard)},
	}
}
