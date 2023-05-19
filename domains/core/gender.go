package core

type Gender int

const (
	// GenderUnknown ...
	GenderUnknown Gender = iota
	// GenderMale ...
	GenderMale
	// GenderFemale ...
	GenderFemale
	// GenderOthers ...
	GenderOthers
)

// String returns a string representation of a gender
func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	case GenderOthers:
		return "others"
	default:
		return "unknown"
	}
}

// GenderFromString return the gender from a given string
func GenderFromString(s string) Gender {
	switch s {
	case GenderFemale.String():
		return GenderFemale
	case GenderMale.String():
		return GenderMale
	case GenderOthers.String():
		return GenderOthers
	default:
		return GenderUnknown
	}
}
