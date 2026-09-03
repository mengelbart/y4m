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
	Presentation      string
	TemporalSampling  string
	ChromaSubsampling string
	Metadata          []string
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
			h.Presentation = string(field[1])
			h.TemporalSampling = string(field[2])
			h.ChromaSubsampling = string(field[3])
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
