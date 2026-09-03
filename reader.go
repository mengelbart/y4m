package y4m

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

var errFrameTooLarge = errors.New("frame too large")

// Option configures a Reader.
type Option func(*Reader)

// WithMaxFrameSize limits the size of the frame buffer a Reader allocates. A
// value of zero or less means no limit.
func WithMaxFrameSize(size int) Option {
	return func(r *Reader) {
		r.maxFrameSize = size
	}
}

type Reader struct {
	reader       *bufio.Reader
	streamHeader *StreamHeader
	maxFrameSize int
}

func NewReader(input io.Reader, opts ...Option) (*Reader, *StreamHeader, error) {
	r := &Reader{
		reader:       bufio.NewReader(input),
		streamHeader: &StreamHeader{},
	}
	for _, opt := range opts {
		opt(r)
	}
	if err := r.streamHeader.parse(r.reader); err != nil {
		return nil, nil, err
	}
	return r, r.streamHeader, nil
}

func (r *Reader) ReadNextFrame() ([]byte, *FrameHeader, error) {
	header := &FrameHeader{}
	if err := header.parse(r.reader, r.streamHeader.Interlacing == ITMixed); err != nil {
		return nil, nil, err
	}
	size, err := r.streamHeader.frameSize()
	if err != nil {
		return nil, nil, err
	}
	if r.maxFrameSize > 0 && size > r.maxFrameSize {
		return nil, nil, fmt.Errorf("%w: %v bytes exceeds maximum of %v", errFrameTooLarge, size, r.maxFrameSize)
	}
	buf := make([]byte, size)
	_, err = io.ReadFull(r.reader, buf)
	if err != nil {
		return nil, nil, err
	}
	return buf, header, nil
}
