package y4m

type ChromaSubsamplingType string

const (
	CST420      ChromaSubsamplingType = "420"
	CST420jpeg  ChromaSubsamplingType = "420jpeg"
	CST420mpeg2 ChromaSubsamplingType = "420mpeg2"
	CST420paldv ChromaSubsamplingType = "420paldv"
	CST411      ChromaSubsamplingType = "411"
	CST422      ChromaSubsamplingType = "422"
	CST444      ChromaSubsamplingType = "444"
	CST444Alpha ChromaSubsamplingType = "444alpha"
	CSTMono     ChromaSubsamplingType = "mono"
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
