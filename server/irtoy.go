package server

import "io"
// import "log"

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
  // log.Printf("writing %v", b)
  n, err := t.serial.Write(b)
  if err != nil { panic(err) }
  if n != len(b) { panic("short write") }
}

func (t *IrToy) read(b []byte) {
  // log.Printf("waiting for %d bytes", len(b))
  n, err := t.serial.Read(b)
  // log.Printf("got %d bytes: %v", n, b)
  if err != nil { panic(err) }
  if n != len(b) { panic("short read") }
}

func (t *IrToy) reset() {
  t.write([]byte{0, 0, 0, 0, 0})
}

func (t *IrToy) transmit(b []byte) {
  var ack [3]byte
  // Enable handshakes, transmit, and whatnot. The acknowledgement is the size
  // of the internal buffer (amount of bytes we can send)
  t.write([]byte{0x26, 0x25, 0x24, 0x03})
  t.read(ack[0:1])
  size := int(ack[0])
  if size < 0 { panic("bad buffer size") }

  // Send all the data in the buffer-size chunks, reading off acknowledgements.
  for i := 0; i < len(b); i += size {
    end := i + size
    if end > len(b) { end = len(b) }
    t.write(b[i:end])
    if end == len(b) {
      t.write([]byte{0xff, 0xff})
    }

    t.read(ack[0:1])
    if int(ack[0]) != size {
      panic("bad acknowledgement")
    }
  }

  // Finally read the transmit count and completion flags
  t.read(ack[0:3])
  if ack[0] != 't' { panic("didn't receive a 't'") }
  if (int(ack[1]) << 8) | int(ack[2]) != len(b) + 2 {
    panic("didn't send all bytes?")
  }
  t.read(ack[0:1])
  if ack[0] != 'C' { panic("didn't actually complete") }
}
