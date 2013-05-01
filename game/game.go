// lots of code from
// https://github.com/BurntSushi/xgbutil/blob/master/_examples/pointer-painting/main.go


// Example pointer-painting shows how to draw on a window, MS Paint style.
// This is an extremely involved example, but it showcases a lot of xgbutil
// and how pieces of it can be tied together.
//
// If you're just starting with xgbutil, I highly recommend checking out the
// other examples before attempting to digest this one.
package game

import (
  "bufio"
  "bytes"
  "flag"
  "fmt"
  "image"
  "math/rand"
  "os"
  "os/exec"
  "time"

  "github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/xgbutil"
  "github.com/BurntSushi/xgbutil/ewmh"
  "github.com/BurntSushi/xgbutil/mousebind"
  "github.com/BurntSushi/xgbutil/xevent"
  "github.com/BurntSushi/xgbutil/xgraphics"
  "github.com/BurntSushi/xgbutil/xwindow"

  "github.com/cgil/verti-gogh/server"
)

const (
  delta = 20.0
  carsize = 30
  width, height = 1100, 600

  calsize        = 40
  CALIBRATE_HIT_THRESH  = 5
  CALIBRATE_DIST_THRESH = 36

  HIT_THRESH = 2500
)

var (
  bg     = xgraphics.BGRA{0xff, 0xff, 0xff, 0xff}
  car    = xgraphics.BGRA{0x00, 0x00, 0x00, 0xff}
  marker = xgraphics.BGRA{0x44, 0x44, 0x44, 0xff}
  green  = xgraphics.BGRA{0x00, 0xff, 0x00, 0xff}
  red    = xgraphics.BGRA{0xff, 0x00, 0x00, 0xff}
)

// Global window variables
var X *xgbutil.XUtil
var win *xwindow.Window
var canvas *xgraphics.Image

var Webcam = flag.Bool("webcam", false, "enable Webcam")

type Dot struct {
  good bool
  found bool
  x int
  y int
}
const XS = 6
const YS = 4

type Game struct {
  dots [XS][YS]Dot
  cmd   *exec.Cmd
  topleft image.Point
  topright image.Point
  botleft image.Point
  botright image.Point

  tr uint8
  tg uint8
  tb uint8
}

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

func (g *Game) atpoint(x, y int) {
  px := width / (XS + 1)
  py := height / (YS + 1)

  xn := (x + px / 2) * (XS + 1) / width
  yn := (y + py / 2) * (YS + 1) / height
  if xn == 0 || yn == 0 || xn > XS || yn > YS { return }

  d := &g.dots[xn - 1][yn - 1]
  dx := d.x - x
  dy := d.y - y
  if !d.found && dx * dx + dy * dy < HIT_THRESH {
    circle(d.x, d.y, 40, bg)
    d.found = true
  }
}

func (g *Game) stop() {
  mousebind.Detach(X, X.RootWin())
  if g.cmd != nil {
    g.cmd.Process.Kill()
    g.cmd.Wait()
    g.cmd = nil
  }
}

func (g *Game) game() {
  // Use the bounds to draw a small dot where we think the black dot on the
  // screen is
  go func() {
    if !*Webcam { return }
    println("tracking", g.tr, g.tg, g.tb)
    pmin := image.Point { X: max(g.topleft.X, g.botleft.X),
                          Y: max(g.topleft.Y, g.topright.Y) }
    pmax := image.Point { X: min(g.topright.X, g.botright.X),
                          Y: min(g.botleft.Y, g.botright.Y) }
    g.cmd = exec.Command("./capture/capture_raw_frames",
                         fmt.Sprintf("0x%02d%02d%02d", g.tr, g.tg, g.tb),
                         fmt.Sprintf("%d", pmin.X),
                         fmt.Sprintf("%d", pmin.Y),
                         fmt.Sprintf("%d", pmax.X),
                         fmt.Sprintf("%d", pmax.Y))
    g.cmd.Stderr = os.Stderr
    out, err := g.cmd.StdoutPipe()
    fatal(err)
    in, err := g.cmd.StdinPipe()
    fatal(err)
    fatal(g.cmd.Start())

    buf := bufio.NewReader(out)
    var p image.Point

    buf.ReadString('\n') // discard first point
    outliers := 0
    for _ = range time.Tick(100 * time.Millisecond) {
      // signal readiness and then wait for it to become available
      in.Write([]byte("go\n"))
      s, err := buf.ReadString('\n')
      if err != nil {
        println("tracker died")
        break
      }

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
      g.atpoint(p.X, p.Y)
    }
  }()

  clearCanvas()

  // Draw a black dot on the cursor
  curx, cury := 0, 0
  xevent.MotionNotifyFun(func(X *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
    ev = compressMotionNotify(X, ev)
    x, y := int(ev.EventX), int(ev.EventY)

    circle(curx, cury, 5, bg)
    circle(x, y, 5, car)
    curx, cury = x, y
    g.atpoint(x, y)
  }).Connect(X, win.Id)

  for x := 0; x < XS; x++ {
    for y := 0; y < YS; y++ {
      g.dots[x][y].found = false
      g.dots[x][y].good = (rand.Int() & 0x1 == 0)

      xloc := width * (x + 1) / (XS + 1)
      yloc := height * (y + 1) / (YS + 1)
      g.dots[x][y].x = xloc
      g.dots[x][y].y = yloc
      if g.dots[x][y].good {
        circle(xloc, yloc, 40, green)
      } else {
        circle(xloc, yloc, 40, red)
      }
    }
  }
}

func locate(c *xgraphics.BGRA) (p image.Point) {
  var buf bytes.Buffer
  cmd := exec.Command("./capture/find_raw",
                      fmt.Sprintf("0x%02x%02x%02x", c.R, c.G, c.B))
  cmd.Stderr = os.Stderr
  cmd.Stdin = nil
  cmd.Stdout = &buf
  err := cmd.Run()
  fatal(err)
  n, err := fmt.Sscanf(buf.String(), "%d %d", &p.X, &p.Y)
  fatal(err)
  if n != 2 { panic("didn't get 2 numbers") }
  return
}

func center() (image.Point, []xgraphics.BGRA) {
  type Bucket struct {
    colors []xgraphics.BGRA
    coord image.Point
  }
  color := xgraphics.BGRA{ A: 0xff }
  buckets := make([]Bucket, 0)

  for {
    // Look through the buckets to see if we have a lot of hits somewhere
    for _, bkt := range buckets {
      if len(bkt.colors) > CALIBRATE_HIT_THRESH {
        circle(width / 2, height / 2, calsize, bg)
        return bkt.coord, bkt.colors
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
        buckets[i].colors = append(buckets[i].colors, color)
        appended = true
        break
      }
    }
    if !appended {
      buckets = append(buckets, Bucket {
        colors: []xgraphics.BGRA{color},
        coord: p,
      })
    }
  }
  panic("not here")
}

func (g *Game) calibrate() {
  if !*Webcam { return }

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
    c, colors := center()

    for _, color := range colors {
      // find the corners
      g.topleft = corner(calsize / 2, calsize / 2, color)
      g.topright = corner(width - calsize / 2, calsize / 2, color)
      g.botleft = corner(calsize / 2, height - calsize / 2, color)
      g.botright = corner(width - calsize / 2, height - calsize / 2, color)

      // validate the corners
      fmt.Printf("center: %v\ntl: %v\ntr: %v\nbl: %v\nbr: %v\n",
                 c, g.topleft, g.topright, g.botleft, g.botright)
      if g.topleft.X > c.X || g.topleft.Y > c.Y {
        println("invalid topleft")
      } else if g.topright.X < c.X || g.topright.Y > c.Y {
        println("invalid topright")
      } else if g.botleft.X > c.X || g.botleft.Y < c.Y {
        println("invalid botleft")
      } else if g.botright.X < c.X || g.botright.Y < c.Y {
        println("invalid botright")
      } else {
        return
      }
      println("make sure the Webcam sees the whole screen")
      time.Sleep(time.Second)
    }

    println("all bad...")
    time.Sleep(2 * time.Second)
  }
}

func Run(c chan server.Packet) {
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

  // calibrate somewhere else and consume this thread for the main loop
  go func() {
    var g Game
    g.tr = 0x1c
    g.tg = 0x1f
    g.tb = 0x24
    for cmd := range c {
      g.stop()
      switch cmd.Cmd {
        case server.Reset:
          g.game()
        case server.Calibrate:
          g.calibrate()
        case server.Target:
          g.tr = cmd.R
          g.tg = cmd.G
          g.tb = cmd.B
          println("requesting", g.tr, g.tg, g.tb)
        case server.Stop:
      }
    }
  }()
  xevent.Main(X)
}
