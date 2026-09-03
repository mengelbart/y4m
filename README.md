# y4m

A Go parser for [YUV4MPEG2](https://wiki.multimedia.cx/index.php/YUV4MPEG2) (`.y4m`) video streams.

Reading only, there is no writer.

## Install

```
go get github.com/mengelbart/y4m
```

## Usage

```go
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mengelbart/y4m"
)

func main() {
	file, err := os.Open("video.y4m")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader, err := y4m.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	header := reader.Header()
	fmt.Println(header.Width, header.Height, header.ChromaSubsampling)

	sizes, err := header.PlaneSizes()
	if err != nil {
		log.Fatal(err)
	}
	for {
		frame, _, err := reader.ReadNextFrame()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		y := frame[:sizes.Y]
		cb := frame[sizes.Y : sizes.Y+sizes.Cb]
		cr := frame[sizes.Y+sizes.Cb : sizes.Y+sizes.Cb+sizes.Cr]
		_, _, _ = y, cb, cr
	}
}
```

`ReadNextFrame` returns the raw planes of one frame, concatenated in the order
Y, Cb, Cr, Alpha. Use `StreamHeader.PlaneSizes` to slice them apart. A clean end
of stream is reported as `io.EOF`, a truncated one as `io.ErrUnexpectedEOF`.

`ReadNextFrameInto` reads into a caller supplied buffer instead, which avoids an
allocation per frame.

`NewReader` takes options, currently `WithMaxFrameSize` to cap the buffer size a
malformed header can ask for.

## Supported chroma subsampling

| tag | plane layout | dimensions |
| --- | --- | --- |
| `420`, `420jpeg`, `420mpeg2`, `420paldv` | Y, Cb, Cr | even width, even height |
| `422` | Y, Cb, Cr | even width |
| `411` | Y, Cb, Cr | width multiple of 4 |
| `444` | Y, Cb, Cr | any |
| `444alpha` | Y, Cb, Cr, Alpha | any |
| `mono` | Y | any |

Headers whose dimensions do not fit the subsampling type are rejected.

## Example

`examples/save-frame-as-png` decodes a single frame to a PNG file:

```
go run ./examples/save-frame-as-png -file video.y4m -skip 30 -png out.png
```
