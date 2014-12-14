package builder

import (
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestLoadDeck(t *testing.T) {
	var loader deckLoader
	loader = &fromDiskDeckLoader{}
	loader.LoadDeck([]string{})
	test.Assert(t, false, "not tested")
}
