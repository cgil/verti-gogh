package game

import (
  "image"
  "log"
  _ "net/http/pprof"

  "github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/xgbutil"
  "github.com/BurntSushi/xgbutil/xevent"
  "github.com/BurntSushi/xgbutil/xgraphics"
)

func fatal(err error) {
  if err != nil {
    log.Panic(err)
  }
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
