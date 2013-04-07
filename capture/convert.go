package main

import "image/jpeg"
import "image/color"
import "image"
import "os"

const WIDTH = 320
const HEIGHT = 240

func main() {
  f, err := os.Open("image.yuv")
  if err != nil { panic(err) }
  defer f.Close()

  im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

  buf := make([]byte, WIDTH * 2)

  for i := 0; i < HEIGHT; i++ {
    n, err := f.Read(buf)
    if err != nil { panic(err) }
    if n != len(buf) { panic("short read") }

    for j := uint(0); j < WIDTH; j++ {
      r, g, b := color.YCbCrToRGB(
        buf[j * 2],
        buf[(j & 0xfffffffe) + 1],
        buf[(j & 0xfffffffe) + 3])
      rgb := color.RGBA {
        R: r,
        G: g,
        B: b,
        A: 1,
      }
      im.SetRGBA(int(j), i, rgb)
    }
  }

  out, err := os.Create("out.jpeg")
  if err != nil { panic(err) }
  defer out.Close()

  jpeg.Encode(out, im, nil)
}
