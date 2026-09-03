package y4m

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var errInvalidFrameHeader = errors.New("invalid frame header")

type FrameHeader struct {
	Presentation     PresentationType
	TemporalSampling TemporalSamplingType
	ChromaSampling   ChromaSamplingType
	Metadata         []string
}

func (h *FrameHeader) parse(reader *bufio.Reader, requireFramingAndSampling bool) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		// A partial line means the stream ended mid-header.
		if errors.Is(err, io.EOF) && len(line) > 0 {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return fmt.Errorf("%w: invalid magic string", errInvalidFrameHeader)
	}
	magicString, fields := fields[0], fields[1:]
	if magicString != "FRAME" {
		return fmt.Errorf("%w: invalid magic string", errInvalidFrameHeader)
	}
	var hasFramingAndSampling bool
	for _, field := range fields {
		switch field[0] {
		case 'I':
			if len(field[1:]) != 3 {
				return fmt.Errorf("%w: invalid framing and sampling tag: %v", errInvalidFrameHeader, field[1:])
			}
			presentation := PresentationType(field[1])
			if !validPresentationTypes[presentation] {
				return fmt.Errorf("%w: invalid presentation type: %v", errInvalidFrameHeader, string(field[1]))
			}
			temporalSampling := TemporalSamplingType(field[2])
			if !validTemporalSamplingTypes[temporalSampling] {
				return fmt.Errorf("%w: invalid temporal sampling type: %v", errInvalidFrameHeader, string(field[2]))
			}
			chromaSampling := ChromaSamplingType(field[3])
			if !validChromaSamplingTypes[chromaSampling] {
				return fmt.Errorf("%w: invalid chroma sampling type: %v", errInvalidFrameHeader, string(field[3]))
			}
			h.Presentation = presentation
			h.TemporalSampling = temporalSampling
			h.ChromaSampling = chromaSampling
			hasFramingAndSampling = true
		case 'X':
			h.Metadata = append(h.Metadata, field[1:])
		default:
			return fmt.Errorf("%w: unknown field: %v", errInvalidFrameHeader, field)
		}
	}
	if requireFramingAndSampling && !hasFramingAndSampling {
		return fmt.Errorf("%w: missing framing and sampling tag", errInvalidFrameHeader)
	}
	return nil
}
