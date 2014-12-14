package builder

import (
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestLoadTemplates(t *testing.T) {
	var loader tmplLoader
	loader = &tmplMapLoader{}
	loader.LoadTemplates([]string{})
	test.Assert(t, false, "not tested")
}
