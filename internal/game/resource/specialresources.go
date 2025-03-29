package resource

// SpecialResource is resource that player need to buildup for happiness
type SpecialResource uint

func (b SpecialResource) String() string {
	switch b {
	case Electricity:
		return "electricity"
	case Waste:
		return "waste"
	default:
		return unknown
	}
}

const (
	Electricity SpecialResource = iota // Main resource for more happier cities
	Waste                              // Is opposite for Electricity, but easier to get
)
