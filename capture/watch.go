package main

import "log"
import "github.com/howeyc/fsnotify"
import "os"
import "os/exec"

func check(err error) {
  if err != nil { panic(err) }
}

func process(file string) {
  bytes, err := exec.Command("./find", file).Output()
  check(err)
  println(string(bytes))
  os.Remove(file)
}

func main() {
  if len(os.Args) > 1 {
    process(os.Args[1])
    return
  }
  path := "./images"
  watcher, err := fsnotify.NewWatcher()
  check(err)
  check(watcher.Watch(path))

  println("watching ", path)
  for {
    select {
      case ev := <-watcher.Event:
        if !ev.IsCreate() { break }
        process(ev.Name)
      case err := <-watcher.Error:
        log.Println("error:", err)
    }
  }

  watcher.Close()
}
