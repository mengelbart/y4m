package y4m

type InterlacingType string

const (
	ITUnknown          InterlacingType = "?"
	ITProgressive      InterlacingType = "p"
	ITTopFieldFirst    InterlacingType = "t"
	ITBottomFieldFirst InterlacingType = "b"
	ITMixed            InterlacingType = "m"
)

var validInterlacingTypes = map[InterlacingType]bool{
	ITUnknown:          true,
	ITProgressive:      true,
	ITTopFieldFirst:    true,
	ITBottomFieldFirst: true,
	ITMixed:            true,
}
