package env

import (
	"github.com/jmorgan1321/SpaceRep/internal/builder"
	"github.com/jmorgan1321/SpaceRep/internal/core"
)

// An Env stores miscellaneous information about our SRS, in order
// to minimize the need for globals and facilitate testing.
type Env struct {
	builder       *builder.Builder
	Deck          *core.Deck
	Distributions core.Distribution
	Seed          int64
	// TmplMap       map[string]*template.Template
	// SessionLength float32 // Duration of this session in minutes
}

type Option func(*Env)

func New(opts ...Option) *Env {
	env := &Env{
		Seed: 42,
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}

// DistributionFunc is used by the environment to determine what
// percentage of the time each bucket should display a test to
// the user.
func DistributionFunc(f func(deck *core.Deck, sessionLength float32) [core.BucketCount]float32) Option {
	return func(e *Env) {
		e.Distributions = f(e.Deck, 10.0)
	}
}

func Deck(d *core.Deck) Option {
	return func(e *Env) {
		e.Deck = d
	}
}

func Seed(seed int64) Option {
	return func(e *Env) {
		e.Seed = seed
	}
}
