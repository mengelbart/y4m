package y4m

import (
	"bufio"
	"io"
)

type Reader struct {
	reader       *bufio.Reader
	streamHeader *StreamHeader
}

func NewReader(input io.Reader) (*Reader, *StreamHeader, error) {
	r := &Reader{
		reader:       bufio.NewReader(input),
		streamHeader: &StreamHeader{},
	}
	if err := r.streamHeader.parse(r.reader); err != nil {
		return nil, nil, err
	}
	return r, r.streamHeader, nil
}

func (r *Reader) ReadNextFrame() ([]byte, *FrameHeader, error) {
	header := &FrameHeader{}
	if r.streamHeader.Interlacing == itMixed {
		header.RequiresFramingAndSampling = true
	}
	if err := header.parse(r.reader); err != nil {
		return nil, nil, err
	}
	size := r.streamHeader.frameSize()
	buf := make([]byte, size)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		return nil, nil, err
	}
	return buf, header, nil
}
