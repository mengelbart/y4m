package y4m

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func parseStreamHeader(input string) (StreamHeader, error) {
	var h StreamHeader
	err := h.parse(bufio.NewReader(strings.NewReader(input)))
	return h, err
}

func TestStreamHeaderParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  StreamHeader
	}{
		{
			name:  "defaults",
			input: "YUV4MPEG2 W16 H16\n",
			want: StreamHeader{
				Width:             16,
				Height:            16,
				ChromaSubsampling: CST420jpeg,
				Interlacing:       ITUnknown,
			},
		},
		{
			name:  "all tags",
			input: "YUV4MPEG2 W1920 H1080 F30000:1001 Ip A1:1 C444 XYSCSS=444\n",
			want: StreamHeader{
				Width:             1920,
				Height:            1080,
				ChromaSubsampling: CST444,
				Interlacing:       ITProgressive,
				FrameRate:         Ratio{Numerator: 30000, Denominator: 1001},
				SampleAspect:      Ratio{Numerator: 1, Denominator: 1},
				Metadata:          []string{"YSCSS=444"},
			},
		},
		{
			name:  "repeated metadata",
			input: "YUV4MPEG2 W16 H16 Xa Xb Xc\n",
			want: StreamHeader{
				Width:             16,
				Height:            16,
				ChromaSubsampling: CST420jpeg,
				Interlacing:       ITUnknown,
				Metadata:          []string{"a", "b", "c"},
			},
		},
		{
			name:  "trailing frame data is not consumed",
			input: "YUV4MPEG2 W16 H16\nFRAME\n",
			want: StreamHeader{
				Width:             16,
				Height:            16,
				ChromaSubsampling: CST420jpeg,
				Interlacing:       ITUnknown,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseStreamHeader(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestStreamHeaderParseChromaSubsamplingTypes(t *testing.T) {
	for chroma := range validChromaSubsamplingTypes {
		t.Run(string(chroma), func(t *testing.T) {
			got, err := parseStreamHeader("YUV4MPEG2 W16 H16 C" + string(chroma) + "\n")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ChromaSubsampling != chroma {
				t.Errorf("got %v, want %v", got.ChromaSubsampling, chroma)
			}
		})
	}
}

func TestStreamHeaderParseInterlacingTypes(t *testing.T) {
	for interlacing := range validInterlacingTypes {
		t.Run(string(interlacing), func(t *testing.T) {
			got, err := parseStreamHeader("YUV4MPEG2 W16 H16 I" + string(interlacing) + "\n")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Interlacing != interlacing {
				t.Errorf("got %v, want %v", got.Interlacing, interlacing)
			}
		})
	}
}

func TestStreamHeaderParseErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty line", "\n"},
		{"invalid magic string", "YUV4MPEG W16 H16\n"},
		{"unknown tag", "YUV4MPEG2 W16 H16 Z1\n"},
		{"missing width", "YUV4MPEG2 H16\n"},
		{"missing height", "YUV4MPEG2 W16\n"},
		{"width not a number", "YUV4MPEG2 Wx H16\n"},
		{"height not a number", "YUV4MPEG2 W16 Hx\n"},
		{"zero width", "YUV4MPEG2 W0 H16\n"},
		{"negative width", "YUV4MPEG2 W-16 H16\n"},
		{"negative height", "YUV4MPEG2 W16 H-16\n"},
		{"invalid chroma subsampling", "YUV4MPEG2 W16 H16 Cbogus\n"},
		{"empty chroma subsampling", "YUV4MPEG2 W16 H16 C\n"},
		{"invalid interlacing", "YUV4MPEG2 W16 H16 Iz\n"},
		{"invalid frame rate", "YUV4MPEG2 W16 H16 F30\n"},
		{"invalid sample aspect", "YUV4MPEG2 W16 H16 Ax:1\n"},
		{"duplicate width", "YUV4MPEG2 W16 W32 H16\n"},
		{"duplicate height", "YUV4MPEG2 W16 H16 H32\n"},
		{"duplicate chroma subsampling", "YUV4MPEG2 W16 H16 C444 C422\n"},
		{"duplicate interlacing", "YUV4MPEG2 W16 H16 Ip Ip\n"},
		{"duplicate frame rate", "YUV4MPEG2 W16 H16 F25:1 F30:1\n"},
		{"duplicate sample aspect", "YUV4MPEG2 W16 H16 A1:1 A2:1\n"},
		{"odd width for 420", "YUV4MPEG2 W15 H16 C420\n"},
		{"odd height for 420", "YUV4MPEG2 W16 H15 C420\n"},
		{"odd width for 422", "YUV4MPEG2 W15 H16 C422\n"},
		{"width not a multiple of 4 for 411", "YUV4MPEG2 W6 H16 C411\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseStreamHeader(tc.input); !errors.Is(err, errInvalidStreamHeader) {
				t.Fatalf("got error %v, want %v", err, errInvalidStreamHeader)
			}
		})
	}
}

func TestStreamHeaderParseOddDimensionsAllowed(t *testing.T) {
	for _, chroma := range []ChromaSubsamplingType{CST444, CST444Alpha, CSTMono} {
		t.Run(string(chroma), func(t *testing.T) {
			if _, err := parseStreamHeader("YUV4MPEG2 W15 H15 C" + string(chroma) + "\n"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestStreamHeaderParseTruncated(t *testing.T) {
	for _, input := range []string{"", "YUV4MPEG2 W16 H1", "YUV4MPEG2 W16 H16"} {
		t.Run(input, func(t *testing.T) {
			if _, err := parseStreamHeader(input); !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("got error %v, want %v", err, io.ErrUnexpectedEOF)
			}
		})
	}
}
