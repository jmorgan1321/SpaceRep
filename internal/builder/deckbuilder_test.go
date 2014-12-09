package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/jmorgan1321/SpaceRep/displays/basic"
	"github.com/jmorgan1321/SpaceRep/internal/core"
	"github.com/jmorgan1321/SpaceRep/internal/test"
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
	checkDisplay(t, i, exp.Display, act.Display)
	checkInfo(t, i, exp.Info, act.Info)
}

func checkDisplay(t *testing.T, i int, exp, act core.Display) {
	// check card display
	e, ok := exp.(*basic.Display)
	a, ok := act.(*basic.Display)
	test.Assert(t, ok, "card %v: expected (*basic.Display), got %v:", i, a)

	test.ExpectEQ(t, e.Word, a.Word, "card %d: Word match", i)
	test.ExpectEQ(t, e.Desc, a.Desc, "card %d: Desc match", i)
	test.ExpectEQ(t, e.Hint, a.Hint, "card %d: Hint match", i)
	test.ExpectEQ(t, e.Comp, a.Comp, "card %d: Comp match", i)
	test.ExpectEQ(t, e.Image, a.Image, "card %d: Image match", i)
}

func checkInfo(t *testing.T, i int, exp, act core.Info) {
	// check card info
	test.ExpectEQ(t, exp.Set, act.Set, "card %d: Set match", i)
	test.ExpectEQ(t, exp.File, act.File, "card %d: Filenames match", i)
	test.ExpectEQ(t, exp.Type, act.Type, "card %d: Type match", i)
	test.ExpectEQ(t, exp.Count, act.Count, "card %d: Count match", i)
	test.ExpectEQ(t, exp.Bucket, act.Bucket, "card %d: Bucket match", i)
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
	usingTestdir(t, func() {
		deck, err := New(dir(testdir.name), DFE(testDFE)).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)

		testdata := []*core.Card{
			&core.Card{
				Display: &basic.Display{Word: "push", Image: "push.jpg", Desc: "push desc", Hint: "push hint", Comp: "git command"},
				Info:    core.Info{Set: "git", File: "push", Type: 1, Count: 2, Bucket: 0},
			},
			&core.Card{
				Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
			},
			&core.Card{
				Display: &basic.Display{Word: "push", Image: "push.jpg", Desc: "push desc", Hint: "push hint", Comp: "git command"},
				Info:    core.Info{Set: "git", File: "push", Type: 2, Count: 1, Bucket: 1},
			},
			&core.Card{
				Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
			},
			&core.Card{
				Display: &basic.Display{Word: "commit", Image: "commit.jpg", Desc: "commit desc", Hint: "commit hint", Comp: "git command"},
				Info:    core.Info{Set: "git", File: "commit", Type: 1, Count: 0, Bucket: 2},
			},
			&core.Card{
				Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
			},
			&core.Card{
				Display: &basic.Display{Word: "commit", Image: "commit.jpg", Desc: "commit desc", Hint: "commit hint", Comp: "git command"},
				Info:    core.Info{Set: "git", File: "commit", Type: 2, Count: 0, Bucket: 3},
			},
			&core.Card{
				Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
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
	})
}

func TestLoadDeck_Specific(t *testing.T) {
	usingTestdir(t, func() {
		// just load 'ppc' deck
		deck, err := New(dir(testdir.name), Deck("ppc"), DFE(testDFE)).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)

		testdata := []*core.Card{
			&core.Card{
				Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
			},
			&core.Card{
				Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
			},
			&core.Card{
				Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
				Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
			},
			&core.Card{
				Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
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
	})
}

func TestLoadDeck_Error(t *testing.T) {
	usingTestdir(t, func() {
		deck, err := New(dir(testdir.name), Deck("DoesNotExist"), DFE(testDFE)).LoadDeck()
		test.Assert(t, err != nil, "expected error")
		test.Expect(t, nil == deck, "deck should be empty")
	})
}

func TestDisplayRendersTemplatesFromDefaultOrLocalDir(t *testing.T) {
	usingTestdir(t, func() {
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
	})
}

func TestDisplaysAreCreatedWithType(t *testing.T) {
	test.Assert(t, false, "untested")
}

func Test_getDeckInfo(t *testing.T) {
	// set up test
	err := os.Mkdir("tmp", 0770)
	test.Assert(t, err == nil, "Mkdir: %v", err)

	f, err := ioutil.TempFile("tmp", "temp")
	test.Assert(t, err == nil, "TempFile: %v", err)

	// clean up
	defer func() {
		f.Close()
		if err := os.RemoveAll("tmp"); err != nil {
			t.Errorf("remove: %v", err)
		}
	}()

	json := []byte(`
	{
	    "Display": "my-display",
	    "Info": [
	        {
	            "File": "add.",
	            "Type": 1,
	            "Count": 0,
	            "Bucket": 0,
	            "FirstSeen": "0001-01-01T00:00:00Z",
	            "LastSeen": "0001-01-01T00:00:00Z"
	        },
	        {
	            "File": "add.",
	            "Type": 2,
	            "Count": 0,
	            "Bucket": 0,
	            "FirstSeen": "0001-01-01T00:00:00Z",
	            "LastSeen": "0001-01-01T00:00:00Z"
	        }
	    ]
	}
	`)
	_, err = f.Write(json)
	test.Assert(t, err == nil, "Write: %v", err)

	exp := deckInfo{
		Display: "my-display",
		Info: []*core.Info{
			&core.Info{File: "add.", Type: 1, Count: 0, Bucket: 0},
			&core.Info{File: "add.", Type: 2, Count: 0, Bucket: 0},
		},
	}

	di, err := getDeckInfo(f.Name())
	test.Assert(t, err == nil, "processCardInfo: %v", err)
	test.ExpectEQ(t, exp.Display, di.Display, "display not read-in properly")

	test.Assert(t, len(di.Info) == 2, "info not read-in properly")
	checkInfo(t, 0, *exp.Info[0], *di.Info[0])
	checkInfo(t, 1, *exp.Info[1], *di.Info[1])
}

func Test_getDataFromDisk(t *testing.T) {
	usingTestdir(t, func() {
		testdata := []core.Display{
			&basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
			&basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
		}

		path := testdir.name + "/html/decks/ppc/cards"
		data, err := getDataFromDisk(path, "basic")
		test.Assert(t, err == nil, "getDataFromDisk: %v", err)

		test.Assert(t, len(data) == len(testdata), "data length mismatch")
		for i, exp := range testdata {
			checkDisplay(t, i, exp, data[i])
		}
	})
}

func Test_makeCards(t *testing.T) {
	testdata := []struct {
		data []core.Display
		info []*core.Info
		exp  []*core.Card
	}{
		// old and new match
		{
			data: []core.Display{
				&basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
				&basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
			},
			info: []*core.Info{
				&core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
				&core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
			},
			exp: []*core.Card{
				&core.Card{
					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
				},
				&core.Card{
					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
				},
				&core.Card{
					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
				},
				&core.Card{
					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
				},
			},
		},
		// deleted 'branch'
		{
			data: []core.Display{
				&basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
			},
			info: []*core.Info{
				&core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
				&core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
			},
			exp: []*core.Card{
				&core.Card{
					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
				},
				&core.Card{
					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
				},
			},
		},
		// add push.
		{
			data: []core.Display{
				&basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
				&basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
			},
			info: []*core.Info{
				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
			},
			exp: []*core.Card{
				&core.Card{
					Display: &basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "cmp", Type: 1, Count: 0, Bucket: 0},
				},
				&core.Card{
					Display: &basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "cmp", Type: 2, Count: 0, Bucket: 0},
				},
				&core.Card{
					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
				},
				&core.Card{
					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
					Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
				},
			},
		},
	}

	for i, tt := range testdata {
		cards := makeCards("ppc", tt.info, tt.data)
		if len(tt.exp) != len(cards) {
			for _, c := range cards {
				fmt.Println(c)
			}
		}
		test.AssertEQ(t, len(tt.exp), len(cards), fmt.Sprintf("len mismatch, test %v", i))
		for j := range tt.exp {
			checkCard(t, i*j+j, tt.exp[i], cards[i])
		}
	}
}

func Test_updateBuckets(t *testing.T) {
	testdata := []struct {
		in, exp []*core.Card
	}{
		{
			in: []*core.Card{
				&core.Card{Info: core.Info{Count: 8, Bucket: 0}},
				&core.Card{Info: core.Info{Count: 4, Bucket: 1}},
				&core.Card{Info: core.Info{Count: 1, Bucket: 2}},
				&core.Card{Info: core.Info{Count: -1, Bucket: 3}},
			},
			exp: []*core.Card{
				&core.Card{Info: core.Info{Count: 0, Bucket: 1}},
				&core.Card{Info: core.Info{Count: 0, Bucket: 2}},
				&core.Card{Info: core.Info{Count: 1, Bucket: 2}},
				&core.Card{Info: core.Info{Count: 0, Bucket: 2}},
			},
		},
	}

	for i, tt := range testdata {
		act := updateBuckets(tt.in)
		for j := range tt.exp {
			checkInfo(t, i*j+j, tt.exp[j].Info, act[j].Info)
		}
	}
}

func Test_writeDeckInfo(t *testing.T) {
	test.Assert(t, false, "untested")
}

func Test_makeCardsSetsDisplayTmpl(t *testing.T) {
	test.Assert(t, false, "untested")
}

func Test_makeCardsSetsInfoSet(t *testing.T) {
	test.Assert(t, false, "untested")
}

func TestLoadCards_CardsInfoGetsCreatedIfItDoesntExist(t *testing.T) {
	test.Assert(t, false, "untested")
}

func TestCardInterface(t *testing.T) {
	testdata := []struct {
		info []Info
		hldr []DispHolder
		out  []Card
		desc string
	}{
		{
			info: []Info{
				{File: "add.data", Type: 0},
				{File: "add.data", Type: 1},
			},
			hldr: []DispHolder{
				{File: "add.data", Disp: &BasicDisplay{Word: "add"}},
			},
			out: []Card{
				&BasicDisplay{
					Stat: Info{File: "add.data", Type: 0},
					Word: "add",
				},
				&BasicDisplay{
					Stat: Info{File: "add.data", Type: 1},
					Word: "add",
				},
			},
			desc: "basic case (no add, no del)",
		},

		{
			info: []Info{},
			hldr: []DispHolder{
				{File: "add.data", Disp: &BasicDisplay{Word: "add"}},
			},
			out: []Card{
				&BasicDisplay{
					Stat: Info{File: "add.data", Type: 0},
					Word: "add",
				},
				&BasicDisplay{
					Stat: Info{File: "add.data", Type: 1},
					Word: "add",
				},
			},
			desc: "add add.data",
		},

		{
			info: []Info{
				{File: "add.data", Type: 0},
				{File: "add.data", Type: 1},
			},
			hldr: []DispHolder{},
			out:  []Card{},
			desc: "del add.data",
		},
	}

	for i, tt := range testdata {
		cards := makeCards2("ppc", tt.info, tt.hldr)
		test.Assert(t, len(cards) == len(tt.out), "test %v: len mismatch", i)

		for j, exp := range tt.out {
			act := cards[j]
			checkCard2(t, j, exp, act)

			if t.Failed() {
				fmt.Printf("test %v, card %v: %v\n", i, j, tt.desc)
			}
		}
	}
}

// TODO: switch to new interfaces.
// TODO: separate out tests into their correct packages.
// TODO: test core package against mock interfaces, add tests for concrete types.

// TODO: test that type renders output correctly.
func TestCardRender2(t *testing.T) {
	scopeTmplMap := core.ScopeTmplMap{
		"ppc": core.TmplMap{
			"thisdoesx": template.Must(template.New("base").Parse(`
				this does x
				`)),
			"xdoesthat": template.Must(template.New("base").Parse(`
				x does that
				`)),
		},
	}

	testdata := []struct {
		c   Card
		exp string
	}{
		{
			c:   &BasicDisplay{Stat: Info{Set: "ppc"}},
			exp: "this does x",
		},
	}

	for i, tt := range testdata {
		s, err := tt.c.Render(scopeTmplMap)
		test.Assert(t, err == nil, "test %v, render: %v", i, err)
		act := strings.TrimSpace(s)
		test.ExpectEQ(t, tt.exp, act, fmt.Sprintf("test %v", i))
	}
}

func checkCard2(t *testing.T, i int, exp, act Card) {
	checkDisplay2(t, i, exp.(Display), act.(Display))
	checkInfo2(t, i, exp.Info(), act.Info())
}

func checkDisplay2(t *testing.T, i int, exp, act Display) {
	// check card display
	e, ok := exp.(*BasicDisplay)
	a, ok := act.(*BasicDisplay)
	test.Assert(t, ok, "card %v: expected (*basic.Display), got %v:", i, a)

	test.ExpectEQ(t, e.Word, a.Word, "card %d: Word match", i)
	test.ExpectEQ(t, e.Desc, a.Desc, "card %d: Desc match", i)
	test.ExpectEQ(t, e.Hint, a.Hint, "card %d: Hint match", i)
	test.ExpectEQ(t, e.Comp, a.Comp, "card %d: Comp match", i)
	test.ExpectEQ(t, e.Image, a.Image, "card %d: Image match", i)
}

func checkInfo2(t *testing.T, i int, exp, act Info) {
	// check card info
	test.ExpectEQ(t, exp.Set, act.Set, "card %d: Set match", i)
	test.ExpectEQ(t, exp.File, act.File, "card %d: Filenames match", i)
	test.ExpectEQ(t, exp.Type, act.Type, "card %d: Type match", i)
	test.ExpectEQ(t, exp.Count, act.Count, "card %d: Count match", i)
	test.ExpectEQ(t, exp.Bucket, act.Bucket, "card %d: Bucket match", i)
}
