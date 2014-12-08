package factory

import (
	"testing"

	"github.com/jmorgan1321/SpaceRep/displays/basic"
	"github.com/jmorgan1321/SpaceRep/internal/core"
	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestDFE(t *testing.T) {
	tests := []struct {
		in  string
		out core.Display
	}{
		{in: "ppc", out: &basic.Display{}},
	}

	for _, tt := range tests {
		d := DFE(tt.in)
		test.ExpectEQ(t, tt.out.Type(), d.Type(), "types match")
	}
}
