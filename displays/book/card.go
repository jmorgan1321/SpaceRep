/*
A book display is used to represent a books main talking points.

Card Types
    - Book Summary:
        Title : The Power of Habit
        Image : some_image.jpg
        Desc  : Habits have triggers, actions, rewards; take time to form;
                are accomplished with little changes, and are often unconscious.
        Topic : Self-Help

        F: What are the main takeaways from the Self-Help book, The Power of Habit? (image)
        B:  Habits have triggers, actions, rewards; take time to form; are
            accomplished with little changes, and are often unconscious.
*/
package book

import "github.com/jmorgan1321/SpaceRep/internal/core"

type Type int

const (
	invalidType = iota

	summary

	lastType
)

func (t Type) String() string {
	switch t {
	default:
		return "Unknown..."
	case summary:
		return "summary"
	}
}

type Card struct {
	core.Info
	Title, Image string
	Desc         string
	Topic        string
}

// Card interface
func (c *Card) Name() string { return c.Title }
func (c *Card) Type() string { return "book" }
func (c *Card) Clone(i core.Info) core.Card {
	return &Card{
		Info: core.Info{
			File:  i.File,
			S:     i.S,
			Type:  i.Type,
			Count: i.Count,
			B:     i.B,
		},
		Title: c.Title,
		Image: c.Image,
		Desc:  c.Desc,
		Topic: c.Topic,
	}
}

func (c *Card) Tmpl() string {
	return Type(c.Info.Type).String()
}

func CreateCardsFromTemplate(c core.Card, i core.Info) []core.Card {
	return []core.Card{
		c.Clone(core.Info{File: i.File, S: i.S, Type: int(summary)}),
	}
}
