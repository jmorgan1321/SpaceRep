/*
package factory is used to create new displays from a text
string.

All new displays should be added to DFE().
*/
package factory

import (
	"errors"

	"github.com/jmorgan1321/SpaceRep/displays/basic"
	"github.com/jmorgan1321/SpaceRep/displays/book"
	"github.com/jmorgan1321/SpaceRep/internal/core"
)

func DFE(s string) (core.Card, error) {
	switch s {
	case "basic":
		return &basic.Card{}, nil
	case "book":
		return &book.Card{}, nil
	}
	return nil, errors.New("unknown card type passed in: " + s)
}

func MakeCards(c core.Card, i core.Info) ([]core.Card, error) {
	switch c.Type() {
	case "basic":
		return basic.CreateCardsFromTemplate(c, i), nil
	case "book":
		return book.CreateCardsFromTemplate(c, i), nil
	}
	return nil, errors.New("unknown display type passed in: " + c.Type())
}
