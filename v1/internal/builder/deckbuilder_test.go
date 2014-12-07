package builder

import (
	"os"
	"strings"
	"testing"

	"github.com/jmorgan1321/SpaceRep/v1/displays/basic"
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
	"github.com/jmorgan1321/SpaceRep/v1/internal/test"
)

// test only function
func dir(s string) Option {
	return func(b *Builder) {
		b.dir = s
	}
}

func testDFE(s string) core.Display {
	switch s {
	case "ppc", "git":
		return &basic.Display{}
	}
	return nil
}

func checkCard(t *testing.T, i int, exp, act *core.Card) {
	// check card display
	e, ok := exp.Display.(*basic.Display)
	a, ok := act.Display.(*basic.Display)
	test.Assert(t, ok, "card %v: expected (*basic.Display), got %v:", i, a)

	test.ExpectEQ(t, e.Word, a.Word, "card %d: Word match", i)
	test.ExpectEQ(t, e.Desc, a.Desc, "card %d: Desc match", i)
	test.ExpectEQ(t, e.Hint, a.Hint, "card %d: Hint match", i)
	test.ExpectEQ(t, e.Comp, a.Comp, "card %d: Comp match", i)
	test.ExpectEQ(t, e.Image, a.Image, "card %d: Image match", i)

	// check card info
	test.ExpectEQ(t, exp.Info.Set, act.Info.Set, "card %d: Set match", i)
	test.ExpectEQ(t, exp.Info.File, act.Info.File, "card %d: Filenames match", i)
	test.ExpectEQ(t, exp.Info.Type, act.Info.Type, "card %d: Type match", i)
	test.ExpectEQ(t, exp.Info.Count, act.Info.Count, "card %d: Count match", i)
	test.ExpectEQ(t, exp.Info.Bucket, act.Info.Bucket, "card %d: Bucket match", i)
}

func TestDeckBuilderDefaultState(t *testing.T) {
	deck := New()
	test.ExpectEQ(t, "", deck.dir, `dir defaults to ""`)
	test.ExpectEQ(t, "all", deck.deck, `deck defaults to "all"`)
	test.Expect(t, nil == deck.dfe, `dfe defaults to nil`)
}

func TestDeckBuilderOption_Deck(t *testing.T) {
	expected := "new deck"
	deck := New(Deck(expected))
	test.ExpectEQ(t, expected, deck.deck, `deck was changed`)
}

func TestDeckBuilderOption_DisplayFactory(t *testing.T) {
	expected := func(string) core.Display {
		return nil
	}
	deck := New(DFE(expected))
	test.Expect(t, deck.dfe != nil, `dfe was changed`)
}

func TestDeckBuilderNew_MultipleOptionsCanBeSetAtOnce(t *testing.T) {
	expectedDir := "new dir"
	expectedDeck := "new deck"
	deck := New(Deck(expectedDeck), dir(expectedDir))
	test.ExpectEQ(t, "new deck", deck.deck, `deck.deck was changed`)
	test.ExpectEQ(t, "new dir", deck.dir, `deck.dir was changed`)
}

func TestLoadDeck_All(t *testing.T) {
	makedir(t)
	// cleanup
	defer func() {
		if err := os.RemoveAll(testdir.name); err != nil {
			t.Errorf("removedir: %v", err)
		}
	}()

	deck, err := New(dir(testdir.name), DFE(testDFE)).LoadDeck()
	test.Assert(t, err == nil, "unexpected error: %v", err)

	testdata := []*core.Card{
		&core.Card{
			Display: &basic.Display{Word: "push", Image: "push.jpg", Desc: "push desc", Hint: "push hint", Comp: 1},
			Info:    core.Info{Set: "git", File: "push", Type: 1, Count: 2, Bucket: 0},
		},
		&core.Card{
			Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
		},
		&core.Card{
			Display: &basic.Display{Word: "push", Image: "push.jpg", Desc: "push desc", Hint: "push hint", Comp: 1},
			Info:    core.Info{Set: "git", File: "push", Type: 2, Count: 1, Bucket: 1},
		},
		&core.Card{
			Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
		},
		&core.Card{
			Display: &basic.Display{Word: "commit", Image: "commit.jpg", Desc: "commit desc", Hint: "commit hint", Comp: 1},
			Info:    core.Info{Set: "git", File: "commit", Type: 1, Count: 0, Bucket: 2},
		},
		&core.Card{
			Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
		},
		&core.Card{
			Display: &basic.Display{Word: "commit", Image: "commit.jpg", Desc: "commit desc", Hint: "commit hint", Comp: 1},
			Info:    core.Info{Set: "git", File: "commit", Type: 2, Count: 0, Bucket: 3},
		},
		&core.Card{
			Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
		},
	}

	i := 0
	for _, b := range deck {
		for _, c := range b {
			checkCard(t, i, testdata[i], c)
			i++
		}
	}
}

func TestLoadDeck_Specific(t *testing.T) {
	makedir(t)
	// cleanup
	defer func() {
		if err := os.RemoveAll(testdir.name); err != nil {
			t.Errorf("removedir: %v", err)
		}
	}()

	// just load 'ppc' deck
	deck, err := New(dir(testdir.name), Deck("ppc"), DFE(testDFE)).LoadDeck()
	test.Assert(t, err == nil, "unexpected error: %v", err)

	testdata := []*core.Card{
		&core.Card{
			Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
		},
		&core.Card{
			Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
		},
		&core.Card{
			Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
		},
		&core.Card{
			Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: 1},
			Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
		},
	}
	i := 0
	for _, b := range deck {
		for _, c := range b {
			checkCard(t, i, testdata[i], c)
			i++
		}
	}
}

func TestLoadDeck_Error(t *testing.T) {
	makedir(t)
	// cleanup
	defer func() {
		if err := os.RemoveAll(testdir.name); err != nil {
			t.Errorf("removedir: %v", err)
		}
	}()

	deck, err := New(dir(testdir.name), Deck("DoesNotExist"), DFE(testDFE)).LoadDeck()
	test.Assert(t, err != nil, "expected error")
	test.Expect(t, nil == deck, "deck should be empty")
}

func TestDisplayRendersTemplatesFromDefaultOrLocalDir(t *testing.T) {
	makedir(t)
	// cleanup
	defer func() {
		if err := os.RemoveAll(testdir.name); err != nil {
			t.Errorf("removedir: %v", err)
		}
	}()

	ppcdeck, err := New(
		dir(testdir.name),
		Deck("ppc"),
		DFE(testDFE),
	).LoadDeck()
	test.Assert(t, err == nil, "unexpected error loading ppc: %v", err)

	gitdeck, err := New(
		dir(testdir.name),
		Deck("git"),
		DFE(testDFE),
	).LoadDeck()
	test.Assert(t, err == nil, "unexpected error loading git: %v", err)

	defaultCard := ppcdeck[core.Daily][0]
	act, err := defaultCard.Render()
	test.Expect(t, err == nil, "unexpected error: %v", err)
	test.Expect(t, strings.Contains(act, "default xdoesthis"),
		"failed to load default")

	localCard := gitdeck[core.Daily][0]
	act, err = localCard.Render()
	test.Expect(t, err == nil, "unexpected error: %v", err)
	test.Expect(t, strings.Contains(act, "local xdoesthis"),
		"failed to load local")
}

func Test_getDisplayErrorsOnUnknownType(t *testing.T) {
	deck := New(dir(testdir.name), Deck("ppc"))
	d, err := deck.getDisplay("unknown")
	test.Assert(t, err != nil, "expected error, got: %v", err)
	test.Expect(t, d == nil, "expected nil, got: %v", d)
}

func TestDisplaysAreCreatedWithType(t *testing.T) {
	deck := New(dir(testdir.name), Deck("ppc"), DFE(testDFE))
	d, err := deck.getDisplay("ppc")
	test.Assert(t, err == nil, "unexpected error: %v", err)

	d.SetType(1)
	test.ExpectEQ(t, "thisdoesx", d.Type(), "wrong type (1)")

	d.SetType(2)
	test.ExpectEQ(t, "xdoesthis", d.Type(), "wrong type (2)")
}

func ExampleBuilder() {

}
