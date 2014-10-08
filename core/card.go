package core

import (
	"github.com/jmorgan1321/SpaceRep/utils"
)

// TODO: This should either be an interface or in the facts lib.  Core shouldn't
// 		 know about concrete types, because then it's harder to make core
//       extensible.
type Type int

func (t Type) String() string {
	switch t {
	case DescCard:
		return "DescCard"
	case WordCard:
		return "WordCard"
	}
	return "Unknown..."
}

const (
	DescCard Type = iota + 1
	WordCard
)

type Display interface{}
type Info struct {
	File string
	Type Type
	// Count stores how many times the user got the correct result on this card.
	//
	// Depending on the bucket the user has to reach a count of 'N' in-order to
	// move this card into the next highest bucket.  A count below zero causes
	// this card to be moved into a lower bucket.
	//
	// TODO: unexport! Temporarily made visible to driver until loading from file happens.
	Count int
	// TODO: unexport! Temporarily made visible to driver until loading from file happens.
	Bucket Bucket
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

// func (c *Card) Count() int {
// 	return c.Count
// }
