package y4m

type InterlacingType string

const (
	itUnknown          InterlacingType = "?"
	itProgressive      InterlacingType = "p"
	itTopFieldFirst    InterlacingType = "t"
	itBottomFieldFirst InterlacingType = "b"
	itMixed            InterlacingType = "m"
)

var validInterlacingTypes = map[InterlacingType]bool{
	itUnknown:          true,
	itProgressive:      true,
	itTopFieldFirst:    true,
	itBottomFieldFirst: true,
	itMixed:            true,
}
