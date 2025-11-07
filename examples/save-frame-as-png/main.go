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
	image := image.NewYCbCr(image.Rect(0, 0, streamHeader.Width, streamHeader.Height), convertSubsampleRatio(streamHeader.ChromaSubsampling))

	ySize := streamHeader.Width * streamHeader.Height
	uSize := ySize / 4
	image.Y = data[:ySize]
	image.Cb = data[ySize : ySize+uSize]
	image.Cr = data[ySize+uSize:]

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	if err = png.Encode(outFile, image); err != nil {
		log.Fatal(err)
	}
}

func convertSubsampleRatio(s y4m.ChromaSubsamplingType) image.YCbCrSubsampleRatio {
	switch s {
	case y4m.CST411:
		return image.YCbCrSubsampleRatio411
	case y4m.CST420:
		return image.YCbCrSubsampleRatio420
	case y4m.CST420jpeg:
		return image.YCbCrSubsampleRatio420
	case y4m.CST420mpeg2:
		return image.YCbCrSubsampleRatio420
	case y4m.CST420paldv:
		return image.YCbCrSubsampleRatio420
	case y4m.CST422:
		return image.YCbCrSubsampleRatio422
	case y4m.CST444:
		return image.YCbCrSubsampleRatio444
	case y4m.CST444Alpha:
		return image.YCbCrSubsampleRatio444
	default:
		panic(fmt.Sprintf("unexpected y4m.ChromaSubsamplingType: %#v", s))
	}
}
