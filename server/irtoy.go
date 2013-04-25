package main

import "io"

type IrToy struct {
  serial io.ReadWriteCloser
}

func NewToy(s io.ReadWriteCloser) *IrToy {
  t := &IrToy { serial: s }
  t.setSamplingMode()
  return t
}

func (t *IrToy) setSamplingMode() {
  t.reset()
  t.write([]byte{'S'})
  var version [3]byte
  t.read(version[0:3])
  if version[0] != 'S' || version[1] != '0' || version[2] != '1' {
    panic("invalid protocol version")
  }
}

func (t *IrToy) write(b []byte) {
  n, err := t.serial.Write(b)
  if err != nil { panic(err) }
  if n != len(b) { panic("short write") }
}

func (t *IrToy) read(b []byte) {
  n, err := t.serial.Read(b)
  if err != nil { panic(err) }
  if n != len(b) { panic("short read") }
}

func (t *IrToy) reset() {
  t.write([]byte{0, 0, 0, 0, 0})
}

func (t *IrToy) transmit(b []byte) {
  b = append(append(b, 0xff), 0xff)
  t.acknowledge([]byte{0x26, 0x25, 0x24, 0x03})
  t.acknowledge(b)

  var ack [4]byte
  t.read(ack[0:4])
  if ack[0] != 't' { panic("didn't receive a 't'") }
  if (int(ack[2]) << 8) | int(ack[1]) != len(b) { panic("didn't send all bytes?") }
  if ack[3] != 'C' { panic("didn't actually complete") }
}

func (t *IrToy) acknowledge(b []byte) {
  var ack [1]byte
  for i := 0; i < len(b); i += 62 {
    end := i + 62
    if end > len(b) { end = len(b) }
    t.write(b[i:end])
    t.read(ack[0:1])
    if ack[0] != 0x3e {
      panic("back acknowledgement")
    }
  }
}
