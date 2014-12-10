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

// HELPERS

func checkCard(t *testing.T, i int, exp, act core.Card) {
	checkDisplay(t, i, exp.(*basic.Card), act.(*basic.Card))
	// TODO: fix
	checkInfo(t, i, exp.(*basic.Card).Info, act.(*basic.Card).Info)
}

func checkDisplay(t *testing.T, i int, exp, act *basic.Card) {
	test.ExpectEQ(t, exp.Word, act.Word, "card %d: Word match", i)
	test.ExpectEQ(t, exp.Desc, act.Desc, "card %d: Desc match", i)
	test.ExpectEQ(t, exp.Hint, act.Hint, "card %d: Hint match", i)
	test.ExpectEQ(t, exp.Comp, act.Comp, "card %d: Comp match", i)
	test.ExpectEQ(t, exp.Image, act.Image, "card %d: Image match", i)
}

func checkInfo(t *testing.T, i int, exp, act *core.Info) {
	test.ExpectEQ(t, exp.Set, act.Set, "card %d: Set match", i)
	test.ExpectEQ(t, exp.File, act.File, "card %d: Filenames match", i)
	test.ExpectEQ(t, exp.Type, act.Type, "card %d: Type match", i)
	test.ExpectEQ(t, exp.Count, act.Count, "card %d: Count match", i)
	test.ExpectEQ(t, exp.Bucket, act.Bucket, "card %d: Bucket match", i)
}

// TESTS

func TestLoadDeck_Specific(t *testing.T) {
	usingTestdir(t, testdir, func() {
		testdata := []core.Card{
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 1, Count: 7, Bucket: 0},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 2, Count: 3, Bucket: 1},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 1, Count: 1, Bucket: 2},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 2, Count: 0, Bucket: 3},
			},
		}

		// just load 'ppc' deck
		deck, err := New(
			dir(testdir.name),
			Deck("ppc"),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _, c := range b {
				checkCard(t, i, testdata[i], c)
				i++
			}
		}
	})
}

func TestMultipleDecksCanBeLoaded(t *testing.T) {
	usingTestdir(t, testdir, func() {
		testdata := []core.Card{
			&basic.Card{
				Word:  "push",
				Image: "push.jpg",
				Desc:  "push desc",
				Hint:  "push hint",
				Comp:  "git command",
				Info:  &core.Info{Set: "git", File: "push.data", Type: 1, Count: 2, Bucket: 0},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 1, Count: 7, Bucket: 0},
			},
			&basic.Card{
				Word:  "push",
				Image: "push.jpg",
				Desc:  "push desc",
				Hint:  "push hint",
				Comp:  "git command",
				Info:  &core.Info{Set: "git", File: "push.data", Type: 2, Count: 1, Bucket: 1},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 2, Count: 3, Bucket: 1},
			},
			&basic.Card{
				Word:  "commit",
				Image: "commit.jpg",
				Desc:  "commit desc",
				Hint:  "commit hint",
				Comp:  "git command",
				Info:  &core.Info{Set: "git", File: "commit.data", Type: 1, Count: 0, Bucket: 2},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 1, Count: 1, Bucket: 2},
			},
			&basic.Card{
				Word:  "commit",
				Image: "commit.jpg",
				Desc:  "commit desc",
				Hint:  "commit hint",
				Comp:  "git command",
				Info:  &core.Info{Set: "git", File: "commit.data", Type: 2, Count: 0, Bucket: 3},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 2, Count: 0, Bucket: 3},
			},
		}

		deck, err := New(
			dir(testdir.name),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _, c := range b {
				checkCard(t, i, testdata[i], c)
				i++
			}
		}
	})
}

func TestLoadDeck_SubPathDecksMustHaveSpecificNames(t *testing.T) {
	usingTestdir(t, nestedTestdir, func() {
		// just load 'ppc' deck
		_, err := New(
			dir(nestedTestdir.name),
			Deck("ppc"),
		).LoadDeck()
		test.Expect(t, err != nil, "didn't recieve expected error")
	})
}

func TestLoadDeck_DecksCanHaveSubPaths(t *testing.T) {
	usingTestdir(t, nestedTestdir, func() {
		testdata := []core.Card{
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "add..data", Type: 1, Count: 7, Bucket: 0},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "add..data", Type: 2, Count: 3, Bucket: 1},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "branch.data", Type: 1, Count: 1, Bucket: 2},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "branch.data", Type: 2, Count: 0, Bucket: 3},
			},
		}

		// load 'facts/ppc' deck
		deck, err := New(
			dir(nestedTestdir.name),
			Deck("facts/ppc"),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _, c := range b {
				checkCard(t, i, testdata[i], c)
				i++
			}
		}
	})
}

func TestLoadDeck_SubPaths_All(t *testing.T) {
	usingTestdir(t, nestedTestdir, func() {
		testdata := []core.Card{
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "add..data", Type: 1, Count: 7, Bucket: 0},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "add..data", Type: 2, Count: 3, Bucket: 1},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "branch.data", Type: 1, Count: 1, Bucket: 2},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "facts\\ppc", File: "branch.data", Type: 2, Count: 0, Bucket: 3},
			},
		}

		// load 'facts/ppc' deck
		deck, err := New(
			dir(nestedTestdir.name),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _, c := range b {
				checkCard(t, i, testdata[i], c)
				i++
			}
		}
	})
}

func TestLoadDeck_DecksCanBeExcludedFromLoading(t *testing.T) {
	usingTestdir(t, testdir, func() {
		testdata := []core.Card{
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 1, Count: 7, Bucket: 0},
			},
			&basic.Card{
				Word:  "add.",
				Image: "add.jpg",
				Desc:  "add. desc",
				Hint:  "add. hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "add..data", Type: 2, Count: 3, Bucket: 1},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 1, Count: 1, Bucket: 2},
			},
			&basic.Card{
				Word:  "branch",
				Image: "branch.jpg",
				Desc:  "branch desc",
				Hint:  "branch hint",
				Comp:  "PowerPC instruction",
				Info:  &core.Info{Set: "ppc", File: "branch.data", Type: 2, Count: 0, Bucket: 3},
			},
		}

		deck, err := New(
			Exclude("git"),
			dir(testdir.name),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _, c := range b {
				checkCard(t, i, testdata[i], c)
				i++
			}
		}
	})
}
func TestLoadDeck_DirsCanBeExcludedFromLoading(t *testing.T) {
	usingTestdir(t, nestedTestdir, func() {
		deck, err := New(
			Exclude("facts/"),
			dir(testdir.name),
		).LoadDeck()
		test.Assert(t, err == nil, "unexpected error: %v", err)
		// TODO: add length check here.  Deck should keep track of it's count?
		// test.AssertEQ(t, len(testdata), len(deck), "Length mismatch")

		i := 0
		for _, b := range deck {
			for _ = range b {
				i++
			}
		}
		test.Assert(t, i == 0, "No cards loaded")
	})
}

func TestLoadDeck_Error(t *testing.T) {
	usingTestdir(t, testdir, func() {
		deck, err := New(
			dir(testdir.name),
			Deck("DoesNotExist"),
		).LoadDeck()
		test.Assert(t, err != nil, "expected error")
		test.Expect(t, nil == deck, "deck should be empty")
	})
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
	test.Assert(t, err == nil, "getDeckInfo: %v", err)
	test.ExpectEQ(t, exp.Display, di.Display, "display not read-in properly")

	test.Assert(t, len(di.Info) == 2, "info not read-in properly")
	checkInfo(t, 0, exp.Info[0], di.Info[0])
	checkInfo(t, 1, exp.Info[1], di.Info[1])
}

func Test_getCardTemplatesFromDisk(t *testing.T) {
	usingTestdir(t, testdir, func() {
		testdata := []CardHolder{
			{File: "add..data", Card: &basic.Card{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"}},
			{File: "branch.data", Card: &basic.Card{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"}},
		}

		path := testdir.name + "/html/decks/ppc/cards"
		data, err := getCardTemplatesFromDisk(path, "basic")
		test.Assert(t, err == nil, "getDataFromDisk: %v", err)

		test.Assert(t, len(data) == len(testdata), "data length mismatch")
		for i, exp := range testdata {
			checkDisplay(t, i, exp.Card.(*basic.Card), data[i].Card.(*basic.Card))
		}
	})
}

// func Test_makeCards(t *testing.T) {
// 	testdata := []struct {
// 		data []core.Display
// 		info []*core.Info
// 		exp  []*core.Card
// 	}{
// 		// old and new match
// 		{
// 			data: []core.Display{
// 				&basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 				&basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 			},
// 			info: []*core.Info{
// 				&core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
// 				&core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
// 				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
// 				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
// 			},
// 			exp: []*core.Card{
// 				&core.Card{
// 					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
// 				},
// 			},
// 		},
// 		// deleted 'branch'
// 		{
// 			data: []core.Display{
// 				&basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 			},
// 			info: []*core.Info{
// 				&core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
// 				&core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
// 				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
// 				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
// 			},
// 			exp: []*core.Card{
// 				&core.Card{
// 					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "add.", Type: 1, Count: 7, Bucket: 0},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "add.", Image: "add.jpg", Desc: "add. desc", Hint: "add. hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "add.", Type: 2, Count: 3, Bucket: 1},
// 				},
// 			},
// 		},
// 		// add push.
// 		{
// 			data: []core.Display{
// 				&basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
// 				&basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 			},
// 			info: []*core.Info{
// 				&core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
// 				&core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
// 			},
// 			exp: []*core.Card{
// 				&core.Card{
// 					Display: &basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "cmp", Type: 1, Count: 0, Bucket: 0},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "cmp", Image: "cmp.jpg", Desc: "cmp desc", Hint: "cmp hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "cmp", Type: 2, Count: 0, Bucket: 0},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "branch", Type: 1, Count: 1, Bucket: 2},
// 				},
// 				&core.Card{
// 					Display: &basic.Display{Word: "branch", Image: "branch.jpg", Desc: "branch desc", Hint: "branch hint", Comp: "PowerPC instruction"},
// 					Info:    core.Info{Set: "ppc", File: "branch", Type: 2, Count: 0, Bucket: 3},
// 				},
// 			},
// 		},
// 	}

// 	for i, tt := range testdata {
// 		cards := makeCards("ppc", tt.info, tt.data)
// 		if len(tt.exp) != len(cards) {
// 			for _, c := range cards {
// 				fmt.Println(c)
// 			}
// 		}
// 		test.AssertEQ(t, len(tt.exp), len(cards), fmt.Sprintf("len mismatch, test %v", i))
// 		for j := range tt.exp {
// 			checkCard(t, i*j+j, tt.exp[i], cards[i])
// 		}
// 	}
// }

// func Test_updateBuckets(t *testing.T) {
// 	testdata := []struct {
// 		in, exp []*core.Card
// 	}{
// 		{
// 			in: []*core.Card{
// 				&core.Card{Info: core.Info{Count: 8, Bucket: 0}},
// 				&core.Card{Info: core.Info{Count: 4, Bucket: 1}},
// 				&core.Card{Info: core.Info{Count: 1, Bucket: 2}},
// 				&core.Card{Info: core.Info{Count: -1, Bucket: 3}},
// 			},
// 			exp: []*core.Card{
// 				&core.Card{Info: core.Info{Count: 0, Bucket: 1}},
// 				&core.Card{Info: core.Info{Count: 0, Bucket: 2}},
// 				&core.Card{Info: core.Info{Count: 1, Bucket: 2}},
// 				&core.Card{Info: core.Info{Count: 0, Bucket: 2}},
// 			},
// 		},
// 	}

// 	for i, tt := range testdata {
// 		act := updateBuckets(tt.in)
// 		for j := range tt.exp {
// 			checkInfo(t, i*j+j, tt.exp[j].Info, act[j].Info)
// 		}
// 	}
// }

// TODO: test these cases

// func TestDisplaysAreCreatedWithType(t *testing.T) {
// 	test.Assert(t, false, "untested")
// }

// func Test_writeDeckInfo(t *testing.T) {
// 	test.Assert(t, false, "untested")
// }

// func Test_makeCardsSetsDisplayTmpl(t *testing.T) {
// 	test.Assert(t, false, "untested")
// }

// func Test_makeCardsSetsInfoSet(t *testing.T) {
// 	test.Assert(t, false, "untested")
// }

// func TestLoadCards_CardsInfoGetsCreatedIfItDoesntExist(t *testing.T) {
// 	test.Assert(t, false, "untested")
// }

// func TestCardInterface(t *testing.T) {
// 	testdata := []struct {
// 		info []Info
// 		hldr []DispHolder
// 		out  []Card
// 		desc string
// 	}{
// 		{
// 			info: []Info{
// 				{File: "add.data", Type: 0},
// 				{File: "add.data", Type: 1},
// 			},
// 			hldr: []DispHolder{
// 				{File: "add.data", Disp: &BasicDisplay{Word: "add"}},
// 			},
// 			out: []Card{
// 				&BasicDisplay{
// 					Stat: Info{File: "add.data", Type: 0},
// 					Word: "add",
// 				},
// 				&BasicDisplay{
// 					Stat: Info{File: "add.data", Type: 1},
// 					Word: "add",
// 				},
// 			},
// 			desc: "basic case (no add, no del)",
// 		},

// 		{
// 			info: []Info{},
// 			hldr: []DispHolder{
// 				{File: "add.data", Disp: &BasicDisplay{Word: "add"}},
// 			},
// 			out: []Card{
// 				&BasicDisplay{
// 					Stat: Info{File: "add.data", Type: 0},
// 					Word: "add",
// 				},
// 				&BasicDisplay{
// 					Stat: Info{File: "add.data", Type: 1},
// 					Word: "add",
// 				},
// 			},
// 			desc: "add add.data",
// 		},

// 		{
// 			info: []Info{
// 				{File: "add.data", Type: 0},
// 				{File: "add.data", Type: 1},
// 			},
// 			hldr: []DispHolder{},
// 			out:  []Card{},
// 			desc: "del add.data",
// 		},
// 	}

// 	for i, tt := range testdata {
// 		cards := makeCards2("ppc", tt.info, tt.hldr)
// 		test.Assert(t, len(cards) == len(tt.out), "test %v: len mismatch", i)

// 		for j, exp := range tt.out {
// 			act := cards[j]
// 			checkCard2(t, j, exp, act)

// 			if t.Failed() {
// 				fmt.Printf("test %v, card %v: %v\n", i, j, tt.desc)
// 			}
// 		}
// 	}
// }

// TODO: separate out tests into their correct packages.
// TODO: test core package against mock interfaces, add tests for concrete types.

// TODO: test that type renders output correctly.
func TestCardRender(t *testing.T) {
	scopeTmplMap := core.ScopeTmplMap{
		"html": core.TmplMap{
			"thisdoesx": template.Must(template.New("base").Parse(`
				default this does x
				`)),
			"xdoesthat": template.Must(template.New("base").Parse(`
				default x does that
				`)),
		},
		"ppc": core.TmplMap{
			"thisdoesx": template.Must(template.New("base").Parse(`
				local this does x
				`)),
		},
	}

	testdata := []struct {
		c   core.Card
		exp string
	}{
		{
			c:   &basic.Card{Info: &core.Info{Set: "ppc", Type: 1}},
			exp: "local this does x",
		},

		{
			c:   &basic.Card{Info: &core.Info{Set: "ppc", Type: 2}},
			exp: "default x does that",
		},
	}

	for i, tt := range testdata {
		s, err := core.Render(scopeTmplMap, tt.c)
		test.Assert(t, err == nil, "test %v, render: %v", i, err)
		act := strings.TrimSpace(s)
		test.ExpectEQ(t, tt.exp, act, fmt.Sprintf("test %v", i))
	}
}
