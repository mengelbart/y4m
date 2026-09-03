package y4m

// InterlacingType is the interlacing of a stream, taken from the stream header
// I tag.
type InterlacingType string

const (
	ITUnknown          InterlacingType = "?" // unknown
	ITProgressive      InterlacingType = "p" // progressive, not interlaced
	ITTopFieldFirst    InterlacingType = "t" // interlaced, top field first
	ITBottomFieldFirst InterlacingType = "b" // interlaced, bottom field first
	ITMixed            InterlacingType = "m" // mixed, every frame header carries its own framing and sampling tag
)

var validInterlacingTypes = map[InterlacingType]bool{
	ITUnknown:          true,
	ITProgressive:      true,
	ITTopFieldFirst:    true,
	ITBottomFieldFirst: true,
	ITMixed:            true,
}
