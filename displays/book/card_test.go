package book

import (
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestDefaultState(t *testing.T) {
	c := Card{}
	test.ExpectEQ(t, "book", c.Type(), "Type() is correct")
}

func TestTypeString(t *testing.T) {
	test.ExpectEQ(t, "Unknown...", invalidType.String(), "invalidType.String() is correct")
	test.ExpectEQ(t, "Unknown...", lastType.String(), "lastType.String() is correct")

	test.ExpectEQ(t, "summary", summary.String(), "summary.String() is correct")
}

func TestClone(t *testing.T) {
}

func TestBucket(t *testing.T) {

}

func TestTmpl(t *testing.T) {
}

func TestStats(t *testing.T) {
}

func TestCreateCardsFromTemplate(t *testing.T) {

}
