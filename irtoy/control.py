import serial
import time
from irtoy import IrToy

# The pre command locks all non-locked cars onto this channel. The pre comand
# will also cease the car's movement. Otherwise all commands last for XXX
# seconds on the car before the car automatically stops that command.
class ShittyChineseCar:
  def __init__(self, toy, channel):
    self.toy = toy
    self.channel = channel

  def pre(self):
    self._send('011011')

  def forward(self):
    self._send('000011')

  def backward(self):
    self._send('110011')

  def left(self):
    self._send('011110')

  def right(self):
    self._send('011000')

  def forwardright(self):
    self._send('000000')

  def forwardleft(self):
    self._send('000110')

  def backwardright(self):
    self._send('110000')

  def backwardleft(self):
    self._send('110110')

  def _send(self, cmd):
    # Bits are each 13 bits: 1CC0XXXXXX0LR (L/R are trim, just be 0)
    # The commands sent to the irtoy have 4 bytes for each bit, where the first
    # pair of bytes is a 16-bit integer representing a high period of activity
    # and the second pair is a low period. All of the low periods are the same,
    # but the high periods vary

    b = [0, 140, 0, 47]             # 1
    if self.channel == 0:           # CC
      b.extend([0, 47, 0, 47])
      b.extend([0, 47, 0, 47])
    elif self.channel == 1:
      b.extend([0, 47, 0, 47])
      b.extend([0, 94, 0, 47])
    elif self.channel == 2:
      b.extend([0, 94, 0, 47])
      b.extend([0, 47, 0, 47])
    else:
      b.extend([0, 94, 0, 47])
      b.extend([0, 94, 0, 47])
    b.extend([0, 47, 0, 47])        # 0

    assert len(cmd) == 6
    for c in cmd:                   # XXXXXX
      if c == '0':
        b.extend([0, 47, 0, 47])
      else:
        b.extend([0, 94, 0, 47])

    b.extend([0, 47, 0, 47])        # 0
    b.extend([0, 47, 0, 47])        # L
    b.extend([0, 47, 0, 47])        # R
    assert len(b) == 13 * 4
    self.toy.transmit(b)


with serial.Serial('/dev/ttyACM0') as serialDevice:
  t = IrToy(serialDevice)
  c = ShittyChineseCar(t, 0)

  print 'forward'
  c.pre()
  c.forward()
  time.sleep(0.2)

  print 'left'
  c.pre()
  c.left()
  time.sleep(0.1)

  print 'forward'
  c.pre()
  c.forward()
  time.sleep(0.1)

