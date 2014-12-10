package builder

import "strings"

type Builder struct {
	dir  string   // if dir != "" then 'html' dir is located (from root) in dir
	deck string   // if != "all" then load try to load that specific deck
	ex   []string // decks to ignore
}

// TODO: unexport
type option func(*Builder)

func New(opts ...option) *Builder {
	bldr := &Builder{
		deck: "all",
	}

	for _, opt := range opts {
		opt(bldr)
	}

	return bldr
}

func Deck(s string) option {
	return func(b *Builder) {
		b.deck = strings.Replace(s, "/", "\\", -1)
	}
}

// TODO: export, make flag for using different dir.
// test only function
func dir(s string) option {
	return func(b *Builder) {
		b.dir = s
	}
}

// Exclude takes a comma separated list of decknames to ignore
// when loading.  If a name ends with a trailing slash, all sub d
// directories will be excluded.
//
// Ex: "git, ppc, facts/" would ignore the git, ppc, and any deck
// under facts.
//
func Exclude(s string) option {
	return func(b *Builder) {
		ex := []string{}
		for _, s := range strings.Split(s, ",") {
			ex = append(ex, strings.TrimSpace(s))
		}
		b.ex = ex
	}
}
