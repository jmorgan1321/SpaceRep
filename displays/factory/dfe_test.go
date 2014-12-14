package factory

import (
	"errors"
	"testing"

	"github.com/jmorgan1321/SpaceRep/displays/basic"
	"github.com/jmorgan1321/SpaceRep/displays/book"
	"github.com/jmorgan1321/SpaceRep/internal/core"
	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestDFE(t *testing.T) {
	tests := []struct {
		in  string
		out core.Card
		err error
	}{
		{in: "unknown", err: errors.New("unknown card type passed in: unknown")},
		{in: "basic", out: &basic.Card{}},
		{in: "book", out: &book.Card{}},
	}

	for i, tt := range tests {
		d, err := DFE(tt.in)
		if tt.err != nil {
			test.ExpectEQ(t, tt.err, err, "test %v: expected error", i)
		} else {
			test.Assert(t, err == nil, "test %v: unexpected error: %v", i, err)
			test.ExpectEQ(t, tt.out.Type(), d.Type(), "test %v: types match", i)
		}
	}
}
