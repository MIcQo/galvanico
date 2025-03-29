package resource

// UserResource is resource which user need to get to make better progress in game
type UserResource uint

func (u UserResource) String() string {
	switch u {
	case Gold:
		return "gold"
	case Ship:
		return "ship"
	case Copper:
		return "copper"
	default:
		return unknown
	}
}

const (
	Gold   UserResource = iota // Gold is base trading currency
	Ship                       // Ship is used for move ships around
	Copper                     // Copper is used as premium resource
)
