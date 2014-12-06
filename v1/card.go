package main

import "errors"

type Deck struct {
	Cards []*Card
}

func NewDeck() *Deck {
	return &Deck{}
}

func (d *Deck) GetNextFlashcard() (*Card, error) {
	if len(d.Cards) == 0 {
		return nil, errors.New("No cards to get")
	}

	return &Card{}, nil
}
func (d *Deck) AddFlashcard(c *Card) {
	d.Cards = append(d.Cards, c)
}

type Card struct{}
