package env

import (
	"testing"
	"time"

	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
	"github.com/jmorgan1321/SpaceRep/v1/internal/test"
)

func TestNew_Default(t *testing.T) {
	env := New()
	test.Assert(t, env != nil, "Env is not nil after New()")
	test.Expect(t, env.Deck == nil, "Empty deck by default")
	test.Expect(t, env.Distributions == [core.BucketCount]float32{}, "empty distributions")
	test.Expect(t, env.Seed == 42, "default seed for PRNG")
}

func TestNew_Deck(t *testing.T) {
	deck := &core.Deck{}
	env := New(Deck(deck))
	test.Assert(t, env.Deck == deck, "Deck was stored")
}

func TestNew_MultipleOptionsCanBeSet(t *testing.T) {
	deck := &core.Deck{}
	env := New(
		Deck(deck),
		DistributionFunc(testStdDistFunc),
	)
	test.Assert(t, env.Deck == deck, "Deck was stored")
	test.Assert(t, env.Distributions != [core.BucketCount]float32{}, "distribution was calculated")
}

func testStdDistFunc(deck *core.Deck, sessionLength float32) [core.BucketCount]float32 {
	return [core.BucketCount]float32{1, 1, 1, 1}
}

func TestNew_DistributionFunc(t *testing.T) {
	env := New(DistributionFunc(testStdDistFunc))
	test.Assert(t, env.Distributions != [core.BucketCount]float32{}, "distribution was calculated")
}

func TestNew_Seed(t *testing.T) {
	env := New(Seed(time.Now().Unix()))
	test.Assert(t, env.Seed != 42, "Seed was set")
}
