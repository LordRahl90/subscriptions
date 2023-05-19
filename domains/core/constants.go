package core

type Swipe int

const (
	// SwipeInvalid invalid swipe
	SwipeInvalid Swipe = iota
	// SwipeYes yes swipe
	SwipeYes
	// SwipeNo no swipe
	SwipeNo
)

// SwipeFromString translate a string to a swipe type
func SwipeFromString(s string) Swipe {
	switch s {
	case SwipeYes.String():
		return SwipeYes
	case SwipeNo.String():
		return SwipeNo
	default:
		return SwipeInvalid
	}
}

func (s Swipe) String() string {
	switch s {
	case SwipeYes:
		return "yes"
	case SwipeNo:
		return "no"
	default:
		return "invalid"
	}
}
