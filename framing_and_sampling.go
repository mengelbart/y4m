package y4m

// PresentationType is the first character of a frame header framing and
// sampling tag.
type PresentationType string

const (
	PTTopFieldFirst          PresentationType = "t" // top field first
	PTTopFieldFirstRepeat    PresentationType = "T" // top field first, repeat the first field
	PTBottomFieldFirst       PresentationType = "b" // bottom field first
	PTBottomFieldFirstRepeat PresentationType = "B" // bottom field first, repeat the first field
	PTProgressiveSingle      PresentationType = "1" // single progressive frame
	PTProgressiveDouble      PresentationType = "2" // progressive frame, repeated once
	PTProgressiveTriple      PresentationType = "3" // progressive frame, repeated twice
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
	TSTProgressive TemporalSamplingType = "p" // both fields sampled at the same time
	TSTInterlaced  TemporalSamplingType = "i" // fields sampled at different times
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
	CHSTFrame   ChromaSamplingType = "p" // chroma subsampled over the whole frame
	CHSTField   ChromaSamplingType = "i" // chroma subsampled over each field separately
	CHSTUnknown ChromaSamplingType = "?" // unknown
)

var validChromaSamplingTypes = map[ChromaSamplingType]bool{
	CHSTFrame:   true,
	CHSTField:   true,
	CHSTUnknown: true,
}
