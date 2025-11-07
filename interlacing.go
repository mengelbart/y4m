package y4m

type InterlacingType string

const (
	itUnknown          InterlacingType = "?"
	itProgressive      InterlacingType = "p"
	itTopFieldFirst    InterlacingType = "t"
	itBottomFieldFirst InterlacingType = "b"
	itMixed            InterlacingType = "m"
)

var validInterlacingTypes = map[InterlacingType]InterlacingType{
	itUnknown:          itUnknown,
	itProgressive:      itProgressive,
	itTopFieldFirst:    itTopFieldFirst,
	itBottomFieldFirst: itBottomFieldFirst,
	itMixed:            itMixed,
}
