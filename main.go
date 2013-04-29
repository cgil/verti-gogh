package main

import "flag"
import "github.com/cgil/verti-gogh/game"
import "github.com/cgil/verti-gogh/server"

func main() {
  flag.Parse()

  c := make(chan server.Command, 0)
  go server.Run(c)
  game.Run(c)
}
