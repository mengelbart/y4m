package y4m

import (
	"errors"
	"testing"
)

func TestRatioParse(t *testing.T) {
	cases := []struct {
		input   string
		want    Ratio
		wantErr bool
	}{
		{input: "1:1", want: Ratio{Numerator: 1, Denominator: 1}},
		{input: "30000:1001", want: Ratio{Numerator: 30000, Denominator: 1001}},
		{input: "0:0", want: Ratio{}},
		{input: "-1:2", want: Ratio{Numerator: -1, Denominator: 2}},
		{input: "", wantErr: true},
		{input: "1", wantErr: true},
		{input: "1:2:3", wantErr: true},
		{input: "a:1", wantErr: true},
		{input: "1:b", wantErr: true},
		{input: ":", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			var got Ratio
			err := got.parse(tc.input)
			if tc.wantErr {
				if !errors.Is(err, errInvalidRatio) {
					t.Fatalf("got error %v, want %v", err, errInvalidRatio)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}
