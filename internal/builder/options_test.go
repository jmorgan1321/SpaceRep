package builder

import (
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestDeckBuilderNew(t *testing.T) {
	deck := New()
	test.ExpectEQ(t, "", deck.dir, `dir defaults to ""`)
	test.ExpectEQ(t, "all", deck.deck, `deck defaults to "all"`)
	test.ExpectEQ(t, 0, len(deck.ex), `excludes nothing by default`)
}

func TestDeckBuilderNew_Deck(t *testing.T) {
	expected := "new deck"
	deck := New(Deck(expected))
	test.ExpectEQ(t, expected, deck.deck, `deck was changed`)
}

func TestDeckBuilderNew_Exclude(t *testing.T) {
	in := "ppc, git,    facts/, lang/dyn/python, lang/dyn/"
	expected := []string{"ppc", "git", "facts/", "lang/dyn/python", "lang/dyn/"}
	deck := New(Exclude(in))
	test.ExpectEQ(t, expected, deck.ex, `exlude list was updated`)
}

func TestDeckBuilderNew_MultipleOptionsCanBeSetAtOnce(t *testing.T) {
	expectedDir := "new dir"
	expectedDeck := "new deck"
	deck := New(Deck(expectedDeck), dir(expectedDir))
	test.ExpectEQ(t, "new deck", deck.deck, `deck.deck was changed`)
	test.ExpectEQ(t, "new dir", deck.dir, `deck.dir was changed`)
}
