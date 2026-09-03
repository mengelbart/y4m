package y4m

// PresentationType is the first character of a frame header framing and
// sampling tag.
type PresentationType string

const (
	PTTopFieldFirst          PresentationType = "t"
	PTTopFieldFirstRepeat    PresentationType = "T"
	PTBottomFieldFirst       PresentationType = "b"
	PTBottomFieldFirstRepeat PresentationType = "B"
	PTProgressiveSingle      PresentationType = "1"
	PTProgressiveDouble      PresentationType = "2"
	PTProgressiveTriple      PresentationType = "3"
)

var validPresentationTypes = map[PresentationType]bool{
	PTTopFieldFirst:          true,
	PTTopFieldFirstRepeat:    true,
	PTBottomFieldFirst:       true,
	PTBottomFieldFirstRepeat: true,
	PTProgressiveSingle:      true,
	PTProgressiveDouble:      true,
	PTProgressiveTriple:      true,
}

// TemporalSamplingType is the second character of a frame header framing and
// sampling tag. It says whether the two fields were sampled at the same time.
type TemporalSamplingType string

const (
	TSTProgressive TemporalSamplingType = "p"
	TSTInterlaced  TemporalSamplingType = "i"
)

var validTemporalSamplingTypes = map[TemporalSamplingType]bool{
	TSTProgressive: true,
	TSTInterlaced:  true,
}

// ChromaSamplingType is the third character of a frame header framing and
// sampling tag. It says whether chroma was subsampled over the whole frame or
// over each field separately.
type ChromaSamplingType string

const (
	CHSTFrame   ChromaSamplingType = "p"
	CHSTField   ChromaSamplingType = "i"
	CHSTUnknown ChromaSamplingType = "?"
)

var validChromaSamplingTypes = map[ChromaSamplingType]bool{
	CHSTFrame:   true,
	CHSTField:   true,
	CHSTUnknown: true,
}
