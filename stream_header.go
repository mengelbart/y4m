package y4m

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var errInvalidStreamHeader = errors.New("invalid stream header")
var errUnsupportedChromaSubsampling = errors.New("unsupported chroma subsampling type")

type StreamHeader struct {
	Width             int
	Height            int
	ChromaSubsampling ChromaSubsamplingType
	Interlacing       InterlacingType
	FrameRate         Ratio
	SampleAspect      Ratio
	Metadata          []string
}

func (h *StreamHeader) frameSize() (int, error) {
	switch h.ChromaSubsampling {
	case CST411, CST420, CST420jpeg, CST420mpeg2, CST420paldv:
		return h.Width * h.Height * 3 / 2, nil

	case CST422:
		return h.Width * h.Height * 2, nil

	case CST444:
		return h.Width * h.Height * 3, nil

	case CST444Alpha:
		return h.Width * h.Height * 4, nil

	case CSTMono:
		return h.Width * h.Height, nil

	default:
		return 0, fmt.Errorf("%w: %v", errUnsupportedChromaSubsampling, h.ChromaSubsampling)
	}
}

func (h *StreamHeader) parse(reader *bufio.Reader) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		// The stream header is mandatory, so any EOF here is truncation.
		if errors.Is(err, io.EOF) {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return fmt.Errorf("%w: too short", errInvalidStreamHeader)
	}
	magicString, fields := fields[0], fields[1:]
	if magicString != "YUV4MPEG2" {
		return fmt.Errorf("%w: invalid magic string", errInvalidStreamHeader)
	}
	// set defaults
	h.ChromaSubsampling = CST420jpeg
	h.Interlacing = itUnknown

	var ok bool
	var foundWidth, foundHeight bool
	for _, field := range fields {
		switch field[0] {
		case 'W':
			h.Width, err = strconv.Atoi(field[1:])
			if err != nil {
				return fmt.Errorf("%w: invalid width: %v", errInvalidStreamHeader, err)
			}
			foundWidth = true
		case 'H':
			h.Height, err = strconv.Atoi(field[1:])
			if err != nil {
				return fmt.Errorf("%w: invalid height: %w", errInvalidStreamHeader, err)
			}
			foundHeight = true
		case 'C':
			h.ChromaSubsampling, ok = validChromaSubsamplingTypes[ChromaSubsamplingType(field[1:])]
			if !ok {
				return fmt.Errorf("%w: invalid chroma subsampling type: %v", errInvalidStreamHeader, field[1:])
			}
		case 'I':
			h.Interlacing, ok = validInterlacingTypes[InterlacingType(field[1:])]
			if !ok {
				return fmt.Errorf("%w: invalid interlacing type: %v", errInvalidStreamHeader, field[1:])
			}
		case 'F':
			if err = h.FrameRate.parse(field[1:]); err != nil {
				return fmt.Errorf("%w: %w", errInvalidStreamHeader, err)
			}
		case 'A':
			if err = h.SampleAspect.parse(field[1:]); err != nil {
				return fmt.Errorf("%w: %w", errInvalidStreamHeader, err)
			}
		case 'X':
			h.Metadata = append(h.Metadata, field[1:])
		default:
			return fmt.Errorf("%w: unknown field: %v", errInvalidStreamHeader, field)
		}
	}
	if !foundWidth || !foundHeight {
		return fmt.Errorf("%w: missing width or heigth", errInvalidStreamHeader)
	}
	if h.Width <= 0 {
		return fmt.Errorf("%w: invalid width: %v", errInvalidStreamHeader, h.Width)
	}
	if h.Height <= 0 {
		return fmt.Errorf("%w: invalid height: %v", errInvalidStreamHeader, h.Height)
	}
	return h.validateDimensions()
}

// validateDimensions checks that the frame dimensions divide evenly into the
// chroma planes of the configured subsampling type:
//
//	420, 420jpeg, 420mpeg2, 420paldv  even width, even height
//	422                               even width
//	411                               width multiple of 4
//	mono, 444, 444alpha               no constraint
func (h *StreamHeader) validateDimensions() error {
	switch h.ChromaSubsampling {
	case CST411:
		if h.Width%4 != 0 {
			return fmt.Errorf("%w: width %v must be a multiple of 4 for chroma subsampling %v", errInvalidStreamHeader, h.Width, h.ChromaSubsampling)
		}

	case CST420, CST420jpeg, CST420mpeg2, CST420paldv:
		if h.Width%2 != 0 || h.Height%2 != 0 {
			return fmt.Errorf("%w: width %v and height %v must be even for chroma subsampling %v", errInvalidStreamHeader, h.Width, h.Height, h.ChromaSubsampling)
		}

	case CST422:
		if h.Width%2 != 0 {
			return fmt.Errorf("%w: width %v must be even for chroma subsampling %v", errInvalidStreamHeader, h.Width, h.ChromaSubsampling)
		}
	}
	return nil
}
