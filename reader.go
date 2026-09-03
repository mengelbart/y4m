package y4m

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
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
	streamHeader StreamHeader
	maxFrameSize int
}

func NewReader(input io.Reader, opts ...Option) (*Reader, error) {
	r := &Reader{
		reader: bufio.NewReader(input),
	}
	for _, opt := range opts {
		opt(r)
	}
	if err := r.streamHeader.parse(r.reader); err != nil {
		return nil, err
	}
	return r, nil
}

// Header returns a copy of the stream header.
func (r *Reader) Header() StreamHeader {
	header := r.streamHeader
	header.Metadata = slices.Clone(r.streamHeader.Metadata)
	return header
}

// ReadNextFrame reads the next frame into a newly allocated buffer.
func (r *Reader) ReadNextFrame() ([]byte, *FrameHeader, error) {
	size, err := r.frameBufferSize()
	if err != nil {
		return nil, nil, err
	}
	buf := make([]byte, size)
	n, header, err := r.ReadNextFrameInto(buf)
	if err != nil {
		return nil, nil, err
	}
	return buf[:n], header, nil
}

// ReadNextFrameInto reads the next frame into buf and returns the number of
// bytes written. It fails with io.ErrShortBuffer, without consuming anything,
// if buf is smaller than a frame.
func (r *Reader) ReadNextFrameInto(buf []byte) (int, *FrameHeader, error) {
	size, err := r.frameBufferSize()
	if err != nil {
		return 0, nil, err
	}
	if len(buf) < size {
		return 0, nil, fmt.Errorf("%w: need %v bytes, got %v", io.ErrShortBuffer, size, len(buf))
	}
	header := &FrameHeader{}
	if err := header.parse(r.reader, r.streamHeader.Interlacing == ITMixed); err != nil {
		return 0, nil, err
	}
	if _, err := io.ReadFull(r.reader, buf[:size]); err != nil {
		// The frame header was consumed, so any EOF here is truncation.
		if errors.Is(err, io.EOF) {
			return 0, nil, io.ErrUnexpectedEOF
		}
		return 0, nil, err
	}
	return size, header, nil
}

// frameBufferSize returns the size of a frame.
func (r *Reader) frameBufferSize() (int, error) {
	size, err := r.streamHeader.frameSize()
	if err != nil {
		return 0, err
	}
	if r.maxFrameSize > 0 && size > r.maxFrameSize {
		return 0, fmt.Errorf("%w: %v bytes exceeds maximum of %v", errFrameTooLarge, size, r.maxFrameSize)
	}
	return size, nil
}
