package main

import (
	"github.com/jmorgan1321/SpaceRep/v1/displays/basic"
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
)

func dfe(s string) core.Display {
	switch s {
	case "ppc", "git":
		return &basic.Display{}
	}
	return nil
}
