package y4m

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

const (
	testStreamHeader = "YUV4MPEG2 W4 H4 F25:1\n"
	testFrameSize    = 4 * 4 * 3 / 2
)

// testStream builds a stream of count frames, each filled with its own index.
func testStream(count int) (string, [][]byte) {
	var stream strings.Builder
	stream.WriteString(testStreamHeader)
	frames := make([][]byte, count)
	for i := range count {
		frames[i] = bytes.Repeat([]byte{byte(i + 1)}, testFrameSize)
		stream.WriteString("FRAME\n")
		stream.Write(frames[i])
	}
	return stream.String(), frames
}

func TestReaderReadNextFrame(t *testing.T) {
	stream, want := testStream(3)
	r, err := NewReader(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, wantFrame := range want {
		got, header, err := r.ReadNextFrame()
		if err != nil {
			t.Fatalf("frame %d: unexpected error: %v", i, err)
		}
		if !bytes.Equal(got, wantFrame) {
			t.Errorf("frame %d: got %v, want %v", i, got, wantFrame)
		}
		if header == nil {
			t.Errorf("frame %d: got nil header", i)
		}
	}
	if _, _, err := r.ReadNextFrame(); !errors.Is(err, io.EOF) {
		t.Fatalf("got error %v, want %v", err, io.EOF)
	}
}

func TestReaderReadNextFrameInto(t *testing.T) {
	stream, want := testStream(3)
	r, err := NewReader(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 4*testFrameSize)
	for i, wantFrame := range want {
		n, _, err := r.ReadNextFrameInto(buf)
		if err != nil {
			t.Fatalf("frame %d: unexpected error: %v", i, err)
		}
		if n != testFrameSize {
			t.Fatalf("frame %d: got %d bytes, want %d", i, n, testFrameSize)
		}
		if !bytes.Equal(buf[:n], wantFrame) {
			t.Errorf("frame %d: got %v, want %v", i, buf[:n], wantFrame)
		}
	}
	if _, _, err := r.ReadNextFrameInto(buf); !errors.Is(err, io.EOF) {
		t.Fatalf("got error %v, want %v", err, io.EOF)
	}
}

func TestReaderReadNextFrameIntoShortBuffer(t *testing.T) {
	stream, want := testStream(1)
	r, err := NewReader(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := r.ReadNextFrameInto(make([]byte, testFrameSize-1)); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("got error %v, want %v", err, io.ErrShortBuffer)
	}
	// The failed call must not have consumed anything.
	got, _, err := r.ReadNextFrame()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, want[0]) {
		t.Errorf("got %v, want %v", got, want[0])
	}
}

func TestReaderTruncated(t *testing.T) {
	stream, _ := testStream(1)
	cases := []struct {
		name  string
		input string
	}{
		{"missing frame data", testStreamHeader + "FRAME\n"},
		{"partial frame data", testStreamHeader + "FRAME\n" + strings.Repeat("\x00", testFrameSize-1)},
		{"partial frame header", stream + "FRAM"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewReader(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var readErr error
			for readErr == nil {
				_, _, readErr = r.ReadNextFrame()
			}
			if !errors.Is(readErr, io.ErrUnexpectedEOF) {
				t.Fatalf("got error %v, want %v", readErr, io.ErrUnexpectedEOF)
			}
		})
	}
}

func TestReaderWithMaxFrameSize(t *testing.T) {
	stream, _ := testStream(1)
	r, err := NewReader(strings.NewReader(stream), WithMaxFrameSize(testFrameSize-1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := r.ReadNextFrame(); !errors.Is(err, errFrameTooLarge) {
		t.Fatalf("got error %v, want %v", err, errFrameTooLarge)
	}

	r, err = NewReader(strings.NewReader(stream), WithMaxFrameSize(testFrameSize))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := r.ReadNextFrame(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReaderHeaderIsCopy(t *testing.T) {
	stream, want := testStream(1)
	r, err := NewReader(strings.NewReader("YUV4MPEG2 W4 H4 Xa\n" + strings.TrimPrefix(stream, testStreamHeader)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	header := r.Header()
	header.Width = 1024
	header.Metadata[0] = "mutated"

	if got := r.Header(); got.Width != 4 || got.Metadata[0] != "a" {
		t.Errorf("reader header was modified: %+v", got)
	}
	got, _, err := r.ReadNextFrame()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, want[0]) {
		t.Errorf("got %v, want %v", got, want[0])
	}
}

func TestReaderMixedInterlacingRequiresFramingAndSampling(t *testing.T) {
	frame := strings.Repeat("\x00", testFrameSize)
	r, err := NewReader(strings.NewReader("YUV4MPEG2 W4 H4 Im\nFRAME\n" + frame))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := r.ReadNextFrame(); !errors.Is(err, errInvalidFrameHeader) {
		t.Fatalf("got error %v, want %v", err, errInvalidFrameHeader)
	}

	r, err = NewReader(strings.NewReader("YUV4MPEG2 W4 H4 Im\nFRAME Itpp\n" + frame))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := r.ReadNextFrame(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewReaderInvalidStreamHeader(t *testing.T) {
	if _, err := NewReader(strings.NewReader("NOPE W4 H4\n")); !errors.Is(err, errInvalidStreamHeader) {
		t.Fatalf("got error %v, want %v", err, errInvalidStreamHeader)
	}
}
