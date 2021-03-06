package server

import "github.com/ziutek/serial"
import "image"
import "image/color"
import "image/jpeg"
import "io"
import "io/ioutil"
import "net/http"
import "os"
import "os/exec"
import "strconv"
import "strings"
import "time"
import _ "net/http/pprof"
import ws "code.google.com/p/go.net/websocket"

const WIDTH = 320
const HEIGHT = 240

type LightningMcQueen struct {
  toy *IrToy
  updown int
  rightleft int
}

type Command int
type Packet struct {
  Cmd Command
  R, G, B uint8
}

const (
  Calibrate Command = iota
  Reset
  Stop
  Target
)

func (car *LightningMcQueen) stop()          { car.send("1132121211212") }
func (car *LightningMcQueen) right()         { car.send("1131122211112") }
func (car *LightningMcQueen) left()          { car.send("1132112211221") }
func (car *LightningMcQueen) forward()       { car.send("1132222211211") }
func (car *LightningMcQueen) backward()      { car.send("1132211211222") }
func (car *LightningMcQueen) fan()           { car.send("1132121212121") }
func (car *LightningMcQueen) forwardright()  { car.send("1131212211122") }
func (car *LightningMcQueen) forwardleft()   { car.send("1131221211111") }
func (car *LightningMcQueen) backwardright() { car.send("1131111212212") }
func (car *LightningMcQueen) backwardleft()  { car.send("1131111211121") }

func (car *LightningMcQueen) send(s string) {
  var cmd [26]byte
  if len(s) != 13 { panic("bad string") }
  for i, c := range s {
    switch c {
      case '1': cmd[2 * i + 1] = 23
      case '2': cmd[2 * i + 1] = 47
      case '3': cmd[2 * i + 1] = 70
      default: panic("bad char: " + string(c))
    }
  }
  car.toy.transmit(cmd[0:26])
}

func Run(c chan Packet, webcam bool) {
  var car *LightningMcQueen

  s, err := serial.Open("/dev/ttyACM0")
  if err != nil {
    println("not actually sending IR commands")
  } else {
    defer s.Close()
    car = &LightningMcQueen { toy: NewToy(s) }
    car.right()
  }

  f, err := os.Open("./server/index2.html")
  if err != nil { panic(err) }
  index, err := ioutil.ReadAll(f)
  if err != nil { panic(err) }
  f.Close()

  f, err = os.Open("./server/control.html")
  if err != nil { panic(err) }
  control, err := ioutil.ReadAll(f)
  if err != nil { panic(err) }
  f.Close()

  f, err = os.Open("./server/image.html")
  if err != nil { panic(err) }
  imhtml, err := ioutil.ReadAll(f)
  if err != nil { panic(err) }
  f.Close()

  event := make(chan string, 0)

  http.Handle("/ws", ws.Handler(func(w *ws.Conn) {
    var message string
    for {
      if ws.Message.Receive(w, &message) != nil { break }
      event <- message
    }
  }))

  srv := http.FileServer(http.Dir("./server/static"))
  http.Handle("/static/", http.StripPrefix("/static", srv))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write(index)
  })
  http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write(control)
  })
  http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write(imhtml)
  })
  http.HandleFunc("/capture.jpg", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    var out io.ReadCloser
    var err error
    if webcam {
      cmd := exec.Command("./capture/yuyv")
      out, err = cmd.StdoutPipe()
      if err != nil { panic(err) }
      cmd.Start()
      defer cmd.Wait()
    } else {
      out, err = os.Open("./capture/image.yuv")
      if err != nil { panic(err) }
    }

    im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    buf := make([]byte, WIDTH * 2)
    for i := 0; i < HEIGHT; i++ {
      n, err := out.Read(buf)
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
    out.Close()
    jpeg.Encode(w, im, nil)
  })

  go func() {
    ticker := time.Tick(400 * time.Millisecond)
    for {
      select {
        case <-ticker:
          if car == nil { continue }
        case s := <-event:
          if car == nil { continue }
          println(s)
          car.stop()
          time.Sleep(25 * time.Millisecond)
          car.stop()
          time.Sleep(25 * time.Millisecond)
          switch s {
            case "-1":  car.updown = 0;  car.rightleft = 0
            // case "0":   car.updown = 0;  car.rightleft = 1
            // case "1":   car.updown = 1;  car.rightleft = 0
            // case "2":   car.updown = 0;  car.rightleft = -1
            // case "3":   car.updown = -1;  car.rightleft = 0
            case "fan": car.fan()
            case "0":  car.updown = 0;  car.rightleft = 1
            case "1":  car.updown = 1;  car.rightleft = 1
            case "2":  car.updown = 1;  car.rightleft = 0
            case "3":  car.updown = 1;  car.rightleft = -1
            case "4":  car.updown = 0;  car.rightleft = -1
            case "5":  car.updown = -1; car.rightleft = -1
            case "6":  car.updown = -1; car.rightleft = 0
            case "7":  car.updown = -1; car.rightleft = 1

            case "calibrate": c <- Packet { Cmd: Calibrate }
            case "reset": c <- Packet { Cmd: Reset }
            case "stop": c <- Packet { Cmd: Stop }

            default:
              p := Packet { Cmd: Target }
              arr := strings.Split(s, ":")
              if arr[0] == "target" && len(arr) == 4 {
                n, err := strconv.ParseUint(arr[1], 10, 8)
                if err != nil { continue }
                p.R = uint8(n)
                n, err = strconv.ParseUint(arr[2], 10, 8)
                if err != nil { continue }
                p.G = uint8(n)
                n, err = strconv.ParseUint(arr[3], 10, 8)
                if err != nil { continue }
                p.B = uint8(n)
                c <- p
              }
          }
      }

      if car.updown == 0 && car.rightleft == 1 {
        car.right()
      } else if car.updown == 0 && car.rightleft == -1 {
        car.left()
      } else if car.updown == 1 && car.rightleft == 0 {
        car.forward()
      } else if car.updown == -1 && car.rightleft == 0 {
        car.backward()
      } else if car.updown == -1 && car.rightleft == -1 {
        car.backwardleft()
      } else if car.updown == 1 && car.rightleft == -1 {
        car.forwardleft()
      } else if car.updown == -1 && car.rightleft == 1 {
        car.backwardright()
      } else if car.updown == 1 && car.rightleft == 1 {
        car.forwardright()
      }
      time.Sleep(25 * time.Millisecond)
      if car.updown == 0 && car.rightleft == 1 {
        car.right()
      } else if car.updown == 0 && car.rightleft == -1 {
        car.left()
      } else if car.updown == 1 && car.rightleft == 0 {
        car.forward()
      } else if car.updown == -1 && car.rightleft == 0 {
        car.backward()
      } else if car.updown == -1 && car.rightleft == -1 {
        car.backwardleft()
      } else if car.updown == 1 && car.rightleft == -1 {
        car.forwardleft()
      } else if car.updown == -1 && car.rightleft == 1 {
        car.backwardright()
      } else if car.updown == 1 && car.rightleft == 1 {
        car.forwardright()
      }
      time.Sleep(25 * time.Millisecond)
    }
  }()

  println("listening")
  http.ListenAndServe(":8000", nil)
}
