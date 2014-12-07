package basic

// CodeComponent is used to differentiate between items, so that cards are
// unique, even with similarly named items.
type CodeComponent int

const (
	Ppc_instr CodeComponent = iota + 1
	Precl
	Emu
	Ppc_reg
)

func (cc CodeComponent) String() string {
	switch cc {
	case Ppc_instr:
		return "PowerPC instruction"
	case Ppc_reg:
		return "PowerPC register"
	case Precl:
		return "Precompiler component"
	case Emu:
		return "Emulator component"
	}
	return "Unknown..."
}
