package main

import (
	"github.com/jmorgan1321/SpaceRep/v1/internal/test"

	"testing"
)

func NewDeckOrAssert(t *testing.T) *Deck {
	d := NewDeck()
	test.Assert(t, d != nil, "expected deck, got: %v", d)
	return d
}

func TestDeckCanBeLoaded(t *testing.T) {

}

func TestDeckCanContainMultipleSets(t *testing.T) {

}

// func TestCantGetCardsFromEmptyDeck(t *testing.T) {
// 	d := NewDeckOrAssert(t)
// 	c, err := d.GetNextFlashcard()
// 	test.Expect(t, err != nil, "expected err, got: %v", err)
// 	test.Expect(t, c == nil, "expected nil, got: %v", c)
// }

// func TestCardsCanBeAddedToDeck(t *testing.T) {
// 	d := NewDeckOrAssert(t)
// 	test.Assert(t, len(d.Cards) == 0, "expected len=0, got: %v", len(d.Cards))

// 	c := &Card{}
// 	d.AddFlashcard(c)

// 	test.Assert(t, len(d.Cards) == 1, "expected len=1, got: %v", len(d.Cards))
// }

// func TestGetFlashcard(t *testing.T) {
// 	d := NewDeckOrAssert(t)

// 	c := &Card{}
// 	d.AddFlashcard(c)

// 	card, err := d.GetNextFlashcard()
// 	test.Assert(t, err == nil, "unexpected error: %v", err)
// 	test.Expect(t, c != nil, "expected card, got: %v", c)
// 	test.ExpectEQ(t, c, card, "Cards should be equal")
// }
