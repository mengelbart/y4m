package y4m

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidStreamHeader = errors.New("invalid stream header")

type StreamHeader struct {
	Width             int
	Height            int
	ChromaSubsampling ChromaSubsamplingType
	Interlacing       InterlacingType
	FrameRate         Ratio
	SampleAspect      Ratio
	Metadata          []string
}

func (h *StreamHeader) frameSize() int {
	switch h.ChromaSubsampling {
	case CST411:
		panic(fmt.Sprintf("not implemented y4m.ChromaSubsamplingType: %#v", h.ChromaSubsampling))

	case CST420, CST420jpeg, CST420mpeg2, CST420paldv:
		return h.Width * h.Height * 3 / 2

	case CST422:
		return h.Width * h.Height * 2

	case CST444:
		return h.Width * h.Height * 3

	case CST444Alpha:
		return h.Width * h.Height * 4

	case CSTMono:
		panic(fmt.Sprintf("not implemented y4m.ChromaSubsamplingType: %#v", h.ChromaSubsampling))

	default:
		panic(fmt.Sprintf("unexpected y4m.ChromaSubsamplingType: %#v", h.ChromaSubsampling))
	}
}

func (h *StreamHeader) parse(reader *bufio.Reader) error {
	line, err := reader.ReadString('\n')
	if err != nil {
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
	return nil
}
