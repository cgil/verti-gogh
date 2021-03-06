package main

import "image/jpeg"
import "image/color"
import "image"
import "os"

const WIDTH = 320
const HEIGHT = 240

func main() {
  f := os.Stdin

  im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

  buf := make([]byte, WIDTH * 2)

  for i := 0; i < HEIGHT; i++ {
    n, err := f.Read(buf)
    if err != nil { panic(err) }
    if n != len(buf) { panic("short read") }

    for j := uint(0); j < WIDTH; j++ {
      yi := j * 2;
      cbi := (j & 0xfffffffe) * 2 + 1
      cri := (j & 0xfffffffe) * 2 + 3

      color := color.YCbCr{ Y: buf[yi], Cb: buf[cbi], Cr: buf[cri] }
      im.Set(int(j), i, color)
    }
  }

  jpeg.Encode(os.Stdout, im, nil)
}
