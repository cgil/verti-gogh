// lots of code from
// https://github.com/BurntSushi/xgbutil/blob/master/_examples/pointer-painting/main.go


// Example pointer-painting shows how to draw on a window, MS Paint style.
// This is an extremely involved example, but it showcases a lot of xgbutil
// and how pieces of it can be tied together.
//
// If you're just starting with xgbutil, I highly recommend checking out the
// other examples before attempting to digest this one.
package main

import (
  "bufio"
  "fmt"
  "image"
  "log"
  "math/rand"
  "os"
  "os/exec"
  "time"
  "net/http"
  "flag"
  _ "net/http/pprof"

  "github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/xgbutil"
  "github.com/BurntSushi/xgbutil/ewmh"
  "github.com/BurntSushi/xgbutil/mousebind"
  "github.com/BurntSushi/xgbutil/xevent"
  "github.com/BurntSushi/xgbutil/xgraphics"
  "github.com/BurntSushi/xgbutil/xwindow"
)

const (
  delta = 20.0
  carsize = 30
  width, height = 1280, 720

  calsize        = 40
  CALIBRATE_HIT_THRESH  = 5
  CALIBRATE_DIST_THRESH = 36

  MONSTER_UPDATE = 200 * time.Millisecond
  MONSTER_CHANGE = 5 * time.Second

  HIT_THRESH = 1600
)

var (
  bg     = xgraphics.BGRA{0xff, 0xff, 0xff, 0xff}
  car    = xgraphics.BGRA{0x00, 0x00, 0x00, 0xff}
  marker = xgraphics.BGRA{0x44, 0x44, 0x44, 0xff}
  green  = xgraphics.BGRA{0x00, 0xff, 0x00, 0xff}
  red    = xgraphics.BGRA{0x00, 0x00, 0xff, 0xff}
)

// Global window variables
var X *xgbutil.XUtil
var win *xwindow.Window
var canvas *xgraphics.Image

var webcam = flag.Bool("webcam", false, "enable webcam")

type Dot struct {
  good bool
  found bool
  x int
  y int
}
const XS = 6
const YS = 4
var dots [XS][YS]Dot

func circle(cx, cy, size int, color xgraphics.BGRA) {
  tipRect := midRect(cx, cy, size, size, width, height)

  // If the rectangle contains no pixels, don't draw anything.
  if tipRect.Empty() {
    return
  }

  tip := canvas.SubImage(tipRect)
  for x := tipRect.Min.X; x < tipRect.Max.X; x++ {
    if x < 0 || x >= width { continue }
    for y := tipRect.Min.Y; y < tipRect.Max.Y; y++ {
      if y < 0 || y >= height { continue }

      dx := x - cx
      dy := y - cy
      if dx * dx + dy * dy < size * size / 4 {
        canvas.SetBGRA(x, y, color)
      } else {
        canvas.SetBGRA(x, y, bg)
      }
    }
  }

  tip.XDraw()
  tip.XPaint(win.Id)
}

func atpoint(x, y int) {
  px := width / (XS + 1)
  py := height / (YS + 1)

  xn := (x + px / 2) * (XS + 1) / width
  yn := (y + py / 2) * (YS + 1) / height
  if xn == 0 || yn == 0 || xn > XS || yn > YS { return }

  d := &dots[xn - 1][yn - 1]
  dx := d.x - x
  dy := d.y - y
  if !d.found && dx * dx + dy * dy < HIT_THRESH {
    circle(d.x, d.y, 40, bg)
    d.found = true
  }
}

func game(topleft, topright, botleft, botright image.Point) {
  // Use the bounds to draw a small dot where we think the black dot on the
  // screen is
  go func() {
    if !*webcam { return }
    pmin := image.Point { X: max(topleft.X, botleft.X),
                          Y: max(topleft.Y, topright.Y) }
    pmax := image.Point { X: min(topright.X, botright.X),
                          Y: min(botleft.Y, botright.Y) }
    cmd := exec.Command("./capture/capture_raw_frames", "0x000000",
                        fmt.Sprintf("%d", pmin.X),
                        fmt.Sprintf("%d", pmin.Y),
                        fmt.Sprintf("%d", pmax.X),
                        fmt.Sprintf("%d", pmax.Y))
    fmt.Printf("%v\n", cmd.Args)
    cmd.Stderr = os.Stderr
    out, err := cmd.StdoutPipe()
    fatal(err)
    in, err := cmd.StdinPipe()
    fatal(err)
    fatal(cmd.Start())

    buf := bufio.NewReader(out)
    var p image.Point

    buf.ReadString('\n') // discard first point
    outliers := 0
    for _ = range time.Tick(100 * time.Millisecond) {
      // signal readiness and then wait for it to become available
      in.Write([]byte("go\n"))
      s, err := buf.ReadString('\n')
      fatal(err)

      var myx, myy int
      n, err := fmt.Sscanf(s, "%d %d", &myy, &myx)
      fatal(err)
      if n != 2 { panic("didn't get 2 ints") }

      myy = (myy - pmin.Y) * height / (pmax.Y - pmin.Y)
      myx = (myx - pmin.X) * width / (pmax.X - pmin.X)
      dx := p.X - myx
      dy := p.Y - myy
      if outliers < 3 && dx * dx + dy * dy > 600 {
        outliers += 1
        continue
      }
      outliers = 0

      circle(p.X, p.Y, 10, bg)
      p.X = myx
      p.Y = myy
      circle(p.X, p.Y, 10, marker)
      atpoint(p.X, p.Y)
    }
  }()

  // Draw a black dot on the cursor
  curx, cury := 0, 0
  xevent.MotionNotifyFun(func(X *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
    ev = compressMotionNotify(X, ev)
    x, y := int(ev.EventX), int(ev.EventY)

    circle(curx, cury, 5, bg)
    circle(x, y, 5, car)
    curx, cury = x, y
    atpoint(x, y)
  }).Connect(X, win.Id)

  for x := 0; x < XS; x++ {
    for y := 0; y < YS; y++ {
      dots[x][y].found = false
      dots[x][y].good = (rand.Int() & 0x1 == 0)

      xloc := width * (x + 1) / (XS + 1)
      yloc := height * (y + 1) / (YS + 1)
      dots[x][y].x = xloc
      dots[x][y].y = yloc
      if dots[x][y].good {
        circle(xloc, yloc, 40, green)
      } else {
        circle(xloc, yloc, 40, red)
      }
    }
  }
}

func locate(c *xgraphics.BGRA) (p image.Point) {
  out, err := exec.Command("./capture/find_raw",
                           fmt.Sprintf("0x%02x%02x%02x", c.R, c.G, c.B)).Output()
  fatal(err)
  n, err := fmt.Sscanf(string(out), "%d %d", &p.X, &p.Y)
  fatal(err)
  if n != 2 { panic("didn't get 2 numbers") }
  return
}

func center() (image.Point, xgraphics.BGRA) {
  type Bucket struct {
    color xgraphics.BGRA
    coord image.Point
    hits  int
  }
  color := xgraphics.BGRA{ A: 0xff }
  buckets := make([]Bucket, 0)

  for {
    fmt.Printf("%v\n", buckets)
    // Look through the buckets to see if we have a lot of hits somewhere
    for _, bkt := range buckets {
      if bkt.hits > CALIBRATE_HIT_THRESH {
        circle(width / 2, height / 2, calsize, bg)
        return bkt.coord, bkt.color
      }
    }

    // Draw a random dot in the center of the screen
    r := rand.Uint32()
    color.R = uint8(r & 0xff)
    color.G = uint8((r >> 8) & 0xff)
    color.B = uint8((r >> 16) & 0xff)
    circle(width / 2, height / 2, calsize, color)
    time.Sleep(200 * time.Millisecond)

    // Find the circle we just drew and append it to the relevant bucket of
    // coordinates
    p := locate(&color)
    appended := false
    for i, _ := range buckets {
      dx := p.X - buckets[i].coord.X
      dy := p.Y - buckets[i].coord.Y
      if dx * dx + dy * dy < CALIBRATE_DIST_THRESH {
        buckets[i].hits++
        appended = true
        break
      }
    }
    if !appended {
      buckets = append(buckets, Bucket {
        color: color,
        coord: p,
        hits: 1,
      })
    }
  }
  panic("not here")
}

func calibrate() {
  var topleft, topright, botleft, botright image.Point
  if !*webcam {
    game(topleft, topright, botleft, botright)
    return
  }

  corner := func(x, y int, color xgraphics.BGRA) image.Point{
    circle(x, y, calsize, color)
    time.Sleep(200 * time.Millisecond)
    ret := locate(&color)
    circle(x, y, calsize, bg)
    return ret
  }

  for {
    // find the center
    clearCanvas()
    c, color := center()

    // find the corners
    topleft = corner(calsize / 2, calsize / 2, color)
    topright = corner(width - calsize / 2, calsize / 2, color)
    botleft = corner(calsize / 2, height - calsize / 2, color)
    botright = corner(width - calsize / 2, height - calsize / 2, color)

    // validate the corners
    fmt.Printf("center: %v\ntl: %v\ntr: %v\nbl: %v\nbr: %v\n",
               c, topleft, topright, botleft, botright)
    if topleft.X > c.X || topleft.Y > c.Y {
      println("invalid topleft")
    } else if topright.X < c.X || topright.Y > c.Y {
      println("invalid topright")
    } else if botleft.X > c.X || botleft.Y < c.Y {
      println("invalid botleft")
    } else if botright.X < c.X || botright.Y < c.Y {
      println("invalid botright")
    } else {
      break
    }

    println("make sure the webcam sees the whole screen")
    time.Sleep(2 * time.Second)
  }

  game(topleft, topright, botleft, botright)
}

func main() {
  flag.Parse()

  var err error
  X, err = xgbutil.NewConn()
  fatal(err)

  mousebind.Initialize(X)
  canvas = xgraphics.New(X, image.Rect(0, 0, width, height))
  canvas.For(func(x, y int) xgraphics.BGRA {
    return bg
  })
  win = canvas.XShowExtra("verti-gogh", true)
  win.Listen(xproto.EventMaskPointerMotion)

  // Atempt to fullscreen
  err = ewmh.WmStateReq(X, win.Id, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")
  fatal(err)

  // profiling
  go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()

  // calibrate somewhere else and consume this thread for the main loop
  go calibrate()
  xevent.Main(X)
}
