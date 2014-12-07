package basic

// Display represent field that we want to display on cards.
type Display struct {
	Word, Image, Desc, Hint string
	Comp                    CodeComponent
	Typ                     Type
}

func (d *Display) Name() string {
	return d.Word
}

func (d *Display) SetType(id int) {
	d.Typ = Type(id)
}

func (d *Display) Type() string {
	return d.Typ.String()
}
