package basic

type Type int

func (t Type) String() string {
	switch t {
	case DescCard:
		return "thisdoesx"
	case WordCard:
		return "xdoesthis"
	}
	return "Unknown..."
}

const (
	DescCard Type = iota + 1
	WordCard
)
