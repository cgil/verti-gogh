package main

import "os"
import "net/http"
import "io/ioutil"
import ws "github.com/kellegous/websocket"

func main() {
  f, err := os.Open("../server/index.html")
  if err != nil { panic(err) }
  index, err := ioutil.ReadAll(f)
  if err != nil { panic(err) }
  f.Close()

  http.Handle("/ws", ws.Handler(func(w *ws.Conn) {
    var message string
    for {
      if ws.Message.Receive(w, &message) != nil {
        break
      }
      println("got", message)
    }
  }))

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write(index)
  })

  http.ListenAndServe(":8000", nil)
}
