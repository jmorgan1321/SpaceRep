/*
package factory is used to create new displays from a text
string.

All new displays should be added to DFE().
*/
package factory

import (
	"errors"

	"github.com/jmorgan1321/SpaceRep/displays/basic"
	"github.com/jmorgan1321/SpaceRep/internal/core"
)

func DFE(s string) (core.Card, error) {
	switch s {
	case "basic":
		return &basic.Card{}, nil
	}
	return nil, errors.New("unknown display type passed in: " + s)
}
