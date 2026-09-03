package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/mengelbart/y4m"
)

func main() {
	path := flag.String("file", "", "Input file")
	out := flag.String("png", "out.png", "Output file")
	skip := flag.Int("skip", 0, "Number of frames to skip")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader, streamHeader, err := y4m.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	for range *skip {
		_, _, err = reader.ReadNextFrame()
		if err != nil {
			log.Fatal(err)
		}
	}
	data, _, err := reader.ReadNextFrame()
	if err != nil {
		log.Fatal(err)
	}
	img, err := frameToImage(streamHeader, data)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	if err = png.Encode(outFile, img); err != nil {
		log.Fatal(err)
	}
}

func frameToImage(h *y4m.StreamHeader, data []byte) (image.Image, error) {
	sizes, err := h.PlaneSizes()
	if err != nil {
		return nil, err
	}
	rect := image.Rect(0, 0, h.Width, h.Height)
	if h.ChromaSubsampling == y4m.CSTMono {
		gray := image.NewGray(rect)
		gray.Pix = data[:sizes.Y]
		return gray, nil
	}
	ratio, err := convertSubsampleRatio(h.ChromaSubsampling)
	if err != nil {
		return nil, err
	}
	cb := sizes.Y
	cr := cb + sizes.Cb
	alpha := cr + sizes.Cr
	img := image.NewYCbCr(rect, ratio)
	img.Y = data[:sizes.Y]
	img.Cb = data[cb:cr]
	img.Cr = data[cr:alpha]
	if sizes.Alpha == 0 {
		return img, nil
	}
	return &image.NYCbCrA{
		YCbCr:   *img,
		A:       data[alpha : alpha+sizes.Alpha],
		AStride: h.Width,
	}, nil
}

func convertSubsampleRatio(s y4m.ChromaSubsamplingType) (image.YCbCrSubsampleRatio, error) {
	switch s {
	case y4m.CST411:
		return image.YCbCrSubsampleRatio411, nil
	case y4m.CST420, y4m.CST420jpeg, y4m.CST420mpeg2, y4m.CST420paldv:
		return image.YCbCrSubsampleRatio420, nil
	case y4m.CST422:
		return image.YCbCrSubsampleRatio422, nil
	case y4m.CST444, y4m.CST444Alpha:
		return image.YCbCrSubsampleRatio444, nil
	default:
		return 0, fmt.Errorf("unsupported chroma subsampling type: %v", s)
	}
}
