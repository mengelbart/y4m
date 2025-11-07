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

var validChromaSubsamplingTypes = map[ChromaSubsamplingType]ChromaSubsamplingType{
	CST420:      CST420,
	CST420jpeg:  CST420jpeg,
	CST420mpeg2: CST420mpeg2,
	CST420paldv: CST420paldv,
	CST411:      CST411,
	CST422:      CST422,
	CST444:      CST444,
	CST444Alpha: CST444Alpha,
	CSTMono:     CSTMono,
}
