package resource

const unknown = "unknown"

// BaseResource is base trade-able unit in game
type BaseResource uint

func (b BaseResource) String() string {
	switch b {
	case Wood:
		return "wood"
	case Water:
		return "water"
	case Iron:
		return "iron"
	case Oil:
		return "oil"
	case Cotton:
		return "cotton"
	default:
		return unknown
	}
}

const (
	Wood   BaseResource = iota // Wood is same unit as Wood
	Water                      // Water is healthier alternative for Wine
	Iron                       // Iron is alternative for Stone
	Oil                        // Oil is an alternative for Crystal
	Cotton                     // Cotton is an alternative for Sulfur
)
