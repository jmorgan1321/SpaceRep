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

// TODO: change the case labels to be names of displays that
//       get stored in a cards.info file.  ie, "basic"
func DFE(s string) (core.Display, error) {
	switch s {
	case "basic":
		return &basic.Display{}, nil
	}
	return nil, errors.New("unknown display type passed in: " + s)
}
