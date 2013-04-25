package main

// import "io/ioutil"
// import "net/http"
// import "os"
// import "time"
// import ws "github.com/kellegous/websocket"
import "github.com/ziutek/serial"

type LightningMcQueen struct {
  toy *IrToy
}

func (car *LightningMcQueen) right() {
  car.send("1131122211112")
}

func (car *LightningMcQueen) send(s string) {
  var cmd [28]byte
  if len(s) != 13 { panic("bad string") }
  for i, c := range s {
    switch c {
      case '1': cmd[2 * i + 1] = 23
      case '2': cmd[2 * i + 1] = 47
      case '3': cmd[2 * i + 1] = 70
      default: panic("bad char: " + string(c))
    }
  }
  cmd[27] = 23
  car.toy.transmit(cmd[0:28])
}

func main() {
  var car *LightningMcQueen

  s, err := serial.Open("/dev/ttyACM0")
  if err != nil {
    println("not actually sending IR commands")
  } else {
    defer s.Close()
    car = &LightningMcQueen { toy: NewToy(s) }
    car.right()
  }

  // f, err := os.Open("./index.html")
  // if err != nil { panic(err) }
  // index, err := ioutil.ReadAll(f)
  // if err != nil { panic(err) }
  // f.Close()
  // left := false
  // right := false
  // forward := false
  // backwards := false
  // fan := false
  // event := make(chan int, 0)

  // http.Handle("/ws", ws.Handler(func(w *ws.Conn) {
  //   var message string
  //   for {
  //     if ws.Message.Receive(w, &message) != nil { break }

  //     println("got", message)

  //     switch message {
  //       case "forward-start": forward = true
  //       case "forward-stop":  forward = false
  //       case "right-start": right = true
  //       case "right-stop":  right = false
  //       case "left-start": left = true
  //       case "left-stop":  left = false
  //       case "backwards-start": backwards = true
  //       case "backwards-stop":  backwards = false
  //       case "fan": fan = !fan
  //     }
  //     event <- 0
  //   }
  // }))

  // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  //   w.WriteHeader(http.StatusOK)
  //   w.Write(index)
  // })

  // go func() {
  //   ticker := time.Tick(100 * time.Millisecond)
  //   for {
  //     select {
  //       case <-ticker:
  //       case <-event:
  //     }
  //     if right {
  //       car.right()
  //     }
  //   }
  // }()

  // println("listening")
  // http.ListenAndServe(":8000", nil)
}
