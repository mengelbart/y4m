package y4m

// ChromaSubsamplingType is the chroma subsampling of a stream, taken from the
// stream header C tag.
type ChromaSubsamplingType string

const (
	CST420      ChromaSubsamplingType = "420"      // 4:2:0 with coincident chroma siting
	CST420jpeg  ChromaSubsamplingType = "420jpeg"  // 4:2:0 with JPEG and MPEG-1 chroma siting
	CST420mpeg2 ChromaSubsamplingType = "420mpeg2" // 4:2:0 with MPEG-2 chroma siting
	CST420paldv ChromaSubsamplingType = "420paldv" // 4:2:0 with PAL-DV chroma siting
	CST411      ChromaSubsamplingType = "411"      // 4:1:1
	CST422      ChromaSubsamplingType = "422"      // 4:2:2
	CST444      ChromaSubsamplingType = "444"      // 4:4:4, chroma not subsampled
	CST444Alpha ChromaSubsamplingType = "444alpha" // 4:4:4 with an alpha plane
	CSTMono     ChromaSubsamplingType = "mono"     // luma only
)

var validChromaSubsamplingTypes = map[ChromaSubsamplingType]bool{
	CST420:      true,
	CST420jpeg:  true,
	CST420mpeg2: true,
	CST420paldv: true,
	CST411:      true,
	CST422:      true,
	CST444:      true,
	CST444Alpha: true,
	CSTMono:     true,
}
