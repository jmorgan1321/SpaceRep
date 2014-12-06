package env

import (
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
	"github.com/jmorgan1321/SpaceRep/v1/internal/test"
	"testing"
)

func TestNew_Default(t *testing.T) {
	env := New()
	test.Assert(t, env != nil, "Env is not nil after New()")
	// test.ExpectEQ(t, 10.0, env.SessionLength, "Session Length defaults to 10mins")
	test.ExpectEQ(t, env.Distributions, [core.BucketCount]float32{}, "empty distributions")
}
