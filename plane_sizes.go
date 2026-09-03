package y4m

import "fmt"

// PlaneSizes holds the size in bytes of each plane of a frame.
type PlaneSizes struct {
	Y     int
	Cb    int
	Cr    int
	Alpha int
}

// Total returns the size of a complete frame in bytes.
func (p PlaneSizes) Total() int {
	return p.Y + p.Cb + p.Cr + p.Alpha
}

// PlaneSizes returns the size of each plane of a frame. The planes are stored
// in a frame in the order Y, Cb, Cr, Alpha.
func (h *StreamHeader) PlaneSizes() (PlaneSizes, error) {
	y := h.Width * h.Height
	switch h.ChromaSubsampling {
	case CST411:
		c := h.Width / 4 * h.Height
		return PlaneSizes{Y: y, Cb: c, Cr: c}, nil

	case CST420, CST420jpeg, CST420mpeg2, CST420paldv:
		c := h.Width / 2 * (h.Height / 2)
		return PlaneSizes{Y: y, Cb: c, Cr: c}, nil

	case CST422:
		c := h.Width / 2 * h.Height
		return PlaneSizes{Y: y, Cb: c, Cr: c}, nil

	case CST444:
		return PlaneSizes{Y: y, Cb: y, Cr: y}, nil

	case CST444Alpha:
		return PlaneSizes{Y: y, Cb: y, Cr: y, Alpha: y}, nil

	case CSTMono:
		return PlaneSizes{Y: y}, nil

	default:
		return PlaneSizes{}, fmt.Errorf("%w: %v", errUnsupportedChromaSubsampling, h.ChromaSubsampling)
	}
}
