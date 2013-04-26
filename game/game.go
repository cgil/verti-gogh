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
  "math"
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
)

const (
  delta = 20.0
  carsize = 30
  width, height = 1024, 500

  calsize        = 40
  CALIBRATE_HIT_THRESH  = 5
  CALIBRATE_DIST_THRESH = 36

  MONSTER_UPDATE = 200 * time.Millisecond
  MONSTER_CHANGE = 5 * time.Second
)

var (
  bg    = xgraphics.BGRA{0xff, 0xff, 0xff, 0xff}
  car   = xgraphics.BGRA{0x00, 0x00, 0x00, 0xff}
  marker = xgraphics.BGRA{0x44, 0x44, 0x44, 0xff}
  green = xgraphics.BGRA{0x00, 0xff, 0x00, 0xff}
  red   = xgraphics.BGRA{0x00, 0x00, 0xff, 0xff}
)

var X *xgbutil.XUtil
var win *xwindow.Window
var canvas *xgraphics.Image

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
        canvas.Set(x, y, color)
      } else {
        canvas.Set(x, y, bg)
      }
    }
  }

  tip.XDraw()
  tip.XPaint(win.Id)
}

func fatal(err error) {
  if err != nil {
    log.Panic(err)
  }
}

func game(topleft, topright, botleft, botright image.Point) {
  // Use the bounds to draw a small dot where we think the black dot on the
  // screen is
  go func() {
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
    for _ = range time.Tick(100 * time.Millisecond) {
      // signal readiness and then wait for it to become available
      in.Write([]byte("go\n"))
      s, err := buf.ReadString('\n')
      fatal(err)

      circle(p.X, p.Y, 10, bg)
      n, err := fmt.Sscanf(s, "%d %d", &p.Y, &p.X)
      p.Y = (p.Y - pmin.Y) * height / (pmax.Y - pmin.Y)
      p.X = (p.X - pmin.X) * width / (pmax.X - pmin.X)
      fatal(err)
      if n != 2 { panic("didn't get 2 ints") }
      circle(p.X, p.Y, 10, marker)
    }
  }()

  // Draw a black dot on the cursor
  curx, cury := 0, 0
  xevent.MotionNotifyFun(func(X *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
    ev = compressMotionNotify(X, ev)
    x, y := int(ev.EventX), int(ev.EventY)

    circle(curx, cury, 40, bg)
    circle(x, y, 40, car)
    curx, cury = x, y
  }).Connect(X, win.Id)

  // Game "monster" loop
  update := time.Tick(MONSTER_UPDATE)
  change := time.Tick(MONSTER_CHANGE)
  x, y := 200, 200
  cur := green
  for {
    select {
      case <-change:
        if cur == green {
          cur = red
        } else {
          cur = green
        }
      case <-update:
    }
    dx := float64(curx - x)
    dy := float64(cury - y)

    // h = <dx, dy>
    // |m * h| == delta
    // sqrt(m * m * dx * dx + m * m * dy * dy) = delta
    // dx * dx + dy * dy = delta * delta
    // m * m = delta * delta / (dx * dx + dy * dy)
    // m = sqrt
    //
    // m * m * delta * delta = dx * dx + dy * dy
    factor := math.Sqrt(delta * delta / (dx * dx + dy * dy))
    if cur == red {
      dx *= factor
      dy *= factor
    } else {
      dx *= -factor
      dy *= -factor
    }
    circle(x, y, carsize, bg)
    x = min(width, max(0, x + int(dx)))
    y = min(height, max(0, y + int(dy)))
    circle(x, y, carsize, cur)
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
  var err error
  X, err = xgbutil.NewConn()
  fatal(err)

  // Whenever the mousebind package is used, you must call Initialize.
  // Similarly for the keybind package.
  mousebind.Initialize(X)

  // Create a new xgraphics.Image. It automatically creates an X pixmap for
  // you, and handles drawing to windows in the XDraw, XPaint and
  // XSurfaceSet functions.
  // N.B. An error is possible since X pixmap allocation can fail.
  canvas = xgraphics.New(X, image.Rect(0, 0, width, height))

  // Color in the background color.
  canvas.For(func(x, y int) xgraphics.BGRA {
    return bg
  })
  win = canvas.XShowExtra("Pointer painting", true)
  win.Listen(xproto.EventMaskPointerMotion)

  // Create a very simple window with dimensions equal to the image.
  // win.Create(im.X.RootWin(), 0, 0, w, h, 0)

  // Make this window close gracefully.
  // win.WMGracefulClose(func(w *xwindow.Window) {
  //   xevent.Detach(w.X, w.Id)
  //   keybind.Detach(w.X, w.Id)
  //   mousebind.Detach(w.X, w.Id)
  //   w.Destroy()

  //   if quit {
  //     xevent.Quit(w.X)
  //   }
  // })

  // Set WM_STATE so it is interpreted as a top-level window.
  // err = icccm.WmStateSet(X, win.Id, &icccm.WmState{
  //   State: icccm.StateNormal,
  // })
  // fatal(err)
  // if err != nil { // not a fatal error
  //   xgbutil.Logger.Printf("Could not set WM_STATE: %s", err)
  // }

  // Set _NET_WM_NAME so it looks nice.
  // err = ewmh.WmNameSet(X, win.Id, "wut")
  // fatal(err)
  // if err != nil { // not a fatal error
  //   xgbutil.Logger.Printf("Could not set _NET_WM_NAME: %s", err)
  // }

  // Paint our image before mapping.
  // im.XSurfaceSet(win.Id)
  // im.XDraw()
  // im.XPaint(win.Id)

  // Now we can map, since we've set all our properties.
  // (The initial map is when the window manager starts managing.)
  // win.Map()

  // Attach event handler for MotionNotify that does not compress events.

  err = ewmh.WmStateReq(X, win.Id, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")
  fatal(err)

  go calibrate()
  xevent.Main(X)
}

// midRect takes an (x, y) position where the pointer was clicked, along with
// the width and height of the thing being drawn and the width and height of
// the canvas, and returns a Rectangle whose midpoint (roughly) is (x, y) and
// whose width and height match the parameters when the rectangle doesn't
// extend past the border of the canvas. Make sure to check if the rectange is
// empty or not before using it!
func midRect(x, y, width, height, canWidth, canHeight int) image.Rectangle {
  return image.Rect(
    x-width/2,   // top left x
    y-height/2, // top left y
    x+width/2,   // bottom right x
    y+height/2, // bottom right y
  )
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func clearCanvas() {
  canvas.For(func(x, y int) xgraphics.BGRA {
    return bg
  })
  canvas.XDraw()
  canvas.XPaint(win.Id)
}

// compressMotionNotify takes a MotionNotify event, and inspects the event
// queue for any future MotionNotify events that can be received without
// blocking. The most recent MotionNotify event is then returned.
// Note that we need to make sure that the Event, Child, Detail, State, Root
// and SameScreen fields are the same to ensure the same window/action is
// generating events. That is, we are only compressing the RootX, RootY,
// EventX and EventY fields.
// This function is not thread safe, since Peek returns a *copy* of the
// event queue---which could be out of date by the time we dequeue events.
func compressMotionNotify(X *xgbutil.XUtil,
ev xevent.MotionNotifyEvent) xevent.MotionNotifyEvent {

  // We force a round trip request so that we make sure to read all
  // available events.
  X.Sync()
  xevent.Read(X, false)

  // The most recent MotionNotify event that we'll end up returning.
  laste := ev

  // Look through each event in the queue. If it's an event and it matches
  // all the fields in 'ev' that are detailed above, then set it to 'laste'.
  // In which case, we'll also dequeue the event, otherwise it will be
  // processed twice!
  // N.B. If our only goal was to find the most recent relevant MotionNotify
  // event, we could traverse the event queue backwards and simply use
  // the first MotionNotify we see. However, this could potentially leave
  // other MotionNotify events in the queue, which we *don't* want to be
  // processed. So we stride along and just pick off MotionNotify events
  // until we don't see any more.
  for i, ee := range xevent.Peek(X) {
    if ee.Err != nil { // This is an error, skip it.
      continue
    }

    // Use type assertion to make sure this is a MotionNotify event.
    if mn, ok := ee.Event.(xproto.MotionNotifyEvent); ok {
      // Now make sure all appropriate fields are equivalent.
      if ev.Event == mn.Event && ev.Child == mn.Child &&
      ev.Detail == mn.Detail && ev.State == mn.State &&
      ev.Root == mn.Root && ev.SameScreen == mn.SameScreen {

        // Set the most recent/valid motion notify event.
        laste = xevent.MotionNotifyEvent{&mn}

        // We cheat and use the stack semantics of defer to dequeue
        // most recent motion notify events first, so that the indices
        // don't become invalid. (If we dequeued oldest first, we'd
        // have to account for all future events shifting to the left
        // by one.)
        defer func(i int) { xevent.DequeueAt(X, i) }(i)
      }
    }
  }

  // This isn't strictly necessary, but is correct. We should update
  // xgbutil's sense of time with the most recent event processed.
  // This is typically done in the main event loop, but since we are
  // subverting the main event loop, we should take care of it.
  X.TimeSet(laste.Time)

  return laste
}
