package y4m

import (
	"errors"
	"testing"
)

func TestPlaneSizes(t *testing.T) {
	cases := []struct {
		chroma ChromaSubsamplingType
		want   PlaneSizes
	}{
		{CST420, PlaneSizes{Y: 256, Cb: 64, Cr: 64}},
		{CST420jpeg, PlaneSizes{Y: 256, Cb: 64, Cr: 64}},
		{CST420mpeg2, PlaneSizes{Y: 256, Cb: 64, Cr: 64}},
		{CST420paldv, PlaneSizes{Y: 256, Cb: 64, Cr: 64}},
		{CST411, PlaneSizes{Y: 256, Cb: 64, Cr: 64}},
		{CST422, PlaneSizes{Y: 256, Cb: 128, Cr: 128}},
		{CST444, PlaneSizes{Y: 256, Cb: 256, Cr: 256}},
		{CST444Alpha, PlaneSizes{Y: 256, Cb: 256, Cr: 256, Alpha: 256}},
		{CSTMono, PlaneSizes{Y: 256}},
	}
	for _, tc := range cases {
		t.Run(string(tc.chroma), func(t *testing.T) {
			h := &StreamHeader{Width: 16, Height: 16, ChromaSubsampling: tc.chroma}
			got, err := h.PlaneSizes()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
			size, err := h.frameSize()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if size != tc.want.Total() {
				t.Errorf("frameSize got %v, want %v", size, tc.want.Total())
			}
		})
	}
}

func TestPlaneSizesUnsupportedChromaSubsampling(t *testing.T) {
	h := &StreamHeader{Width: 16, Height: 16, ChromaSubsampling: "bogus"}
	if _, err := h.PlaneSizes(); !errors.Is(err, errUnsupportedChromaSubsampling) {
		t.Fatalf("got error %v, want %v", err, errUnsupportedChromaSubsampling)
	}
	if _, err := h.frameSize(); !errors.Is(err, errUnsupportedChromaSubsampling) {
		t.Fatalf("got error %v, want %v", err, errUnsupportedChromaSubsampling)
	}
}
