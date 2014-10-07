/*
package work is used to represent work information I want to encapsulate.

Card Types
    - Code Components - generate the following:
        - What does this do?
            Word: GuestFunction
            Img : welcome_guest.jpg
            Desc: stores info about supposed functions from 360 PPC Binaries
            Hint: bob wrote this, drunk
            Comp: precl_comp

            F: "What does the precompiler's GuestFunction do?", image
            B: GuestFunction, precompiler,  stores info..., hint

            Word: bcctr
            Img : tree.jpg
            Desc: branches from the counter register and affects registers x,y,z
            Hint: like that other one
            Comp: ppc_instr

            F: "what does the ppc instr's bcctr do?", image
            B: bbctr, ppc instr, branches from the counter..., hint

        - What does this?
            Word: keInitMemMgr
            Img : some_image.jpg
            Desc: get used by the guest kernel, unknowingly, when it allocates memory.
            Hint: 8GB reserved
            Comp: emu_comp

            F: What emulator component get used by the guest kernel, unknowingly...
            B: keInitMemMgr, emulator, Image, hint
*/
package work

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

// Display represent field that we want to display on cards.
type Display struct {
	Word  string
	Image string
	Desc  string
	Hint  string
	Comp  CodeComponent
}
