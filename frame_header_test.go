package y4m

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func parseFrameHeader(input string, requireFramingAndSampling bool) (FrameHeader, error) {
	var h FrameHeader
	err := h.parse(bufio.NewReader(strings.NewReader(input)), requireFramingAndSampling)
	return h, err
}

func TestFrameHeaderParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  FrameHeader
	}{
		{
			name:  "bare frame",
			input: "FRAME\n",
			want:  FrameHeader{},
		},
		{
			name:  "framing and sampling tag",
			input: "FRAME Itpi\n",
			want: FrameHeader{
				Presentation:     PTTopFieldFirst,
				TemporalSampling: TSTProgressive,
				ChromaSampling:   CHSTField,
			},
		},
		{
			name:  "unknown sampling",
			input: "FRAME I1??\n",
			want: FrameHeader{
				Presentation:     PTProgressiveSingle,
				TemporalSampling: TSTUnknown,
				ChromaSampling:   CHSTUnknown,
			},
		},
		{
			name:  "metadata",
			input: "FRAME Xa Xb\n",
			want:  FrameHeader{Metadata: []string{"a", "b"}},
		},
		{
			name:  "trailing frame data is not consumed",
			input: "FRAME\nsome data",
			want:  FrameHeader{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseFrameHeader(tc.input, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestFrameHeaderParsePresentationTypes(t *testing.T) {
	for presentation := range validPresentationTypes {
		t.Run(string(presentation), func(t *testing.T) {
			got, err := parseFrameHeader("FRAME I"+string(presentation)+"pp\n", true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Presentation != presentation {
				t.Errorf("got %v, want %v", got.Presentation, presentation)
			}
		})
	}
}

func TestFrameHeaderParseErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty line", "\n"},
		{"invalid magic string", "NOTAFRAME\n"},
		{"unknown tag", "FRAME Z1\n"},
		{"tag too short", "FRAME Itp\n"},
		{"tag too long", "FRAME Itppp\n"},
		{"invalid presentation", "FRAME Ixpp\n"},
		{"invalid temporal sampling", "FRAME Itzp\n"},
		{"invalid chroma sampling", "FRAME Itpz\n"},
		{"duplicate framing and sampling tag", "FRAME Itpp Itpp\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseFrameHeader(tc.input, false); !errors.Is(err, errInvalidFrameHeader) {
				t.Fatalf("got error %v, want %v", err, errInvalidFrameHeader)
			}
		})
	}
}

func TestFrameHeaderParseRequiresFramingAndSampling(t *testing.T) {
	if _, err := parseFrameHeader("FRAME\n", true); !errors.Is(err, errInvalidFrameHeader) {
		t.Fatalf("got error %v, want %v", err, errInvalidFrameHeader)
	}
	if _, err := parseFrameHeader("FRAME Itpp\n", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFrameHeaderParseEOF(t *testing.T) {
	if _, err := parseFrameHeader("", false); !errors.Is(err, io.EOF) {
		t.Fatalf("got error %v, want %v", err, io.EOF)
	}
	if _, err := parseFrameHeader("FRAM", false); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("got error %v, want %v", err, io.ErrUnexpectedEOF)
	}
}
