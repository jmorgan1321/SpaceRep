package env

import (
	"github.com/jmorgan1321/srs/core"
	"text/template"
)

// An Env stores miscellaneous information about our SRS, inorder
// to minimize the need for globals and facilitate testing.
type Env struct {
	SessionLength float32 // Duration of this session in minutes
	Cards         [core.BucketCount][]*core.Card
	Distributions [core.BucketCount]float32
	TmplMap       map[string]*template.Template
}

func New() *Env {
	env := &Env{
		SessionLength: 10.0,
		Cards:         LoadCardsFunc(),
		TmplMap:       map[string]*template.Template{},
	}
	env.Distributions = DistributionFunc(env.Cards, env.SessionLength)

	return env
}

// TODO: remove when we have serialization.
var LoadCardsFunc func() [core.BucketCount][]*core.Card

// DistributionFunc is used by the environment to determine what
// percentage of the time each bucket should display a test to
// the user.
var DistributionFunc func(cards [core.BucketCount][]*core.Card, sessionLength float32) [core.BucketCount]float32
