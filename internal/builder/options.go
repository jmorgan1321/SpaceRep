package builder

import "strings"

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
// Ex: "git,ppc,facts/" would ignore the git, ppc, and any deck
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
