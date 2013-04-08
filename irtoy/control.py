import serial
import time
from irtoy import IrToy


print "opening"
with serial.Serial('/dev/ttyACM0') as serialDevice:
  print "opened"
  t = IrToy(serialDevice)
  print "created"

  # print 'receiving'
  # code = t.receive()
  # print code
  print "sending\n\n\n"
  # t.transmit([0, 140, 0, 47,
  #             0, 47, 0, 47,
  #             0, 94, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47,
  #             0, 94, 0, 47,
  #             0, 94, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47,
  #             0, 47, 0, 47])
  # '101000001100# 0'
  # print t.receive()

  pre  =     [0, 140, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 8, 47,
              ]
  foo  =     [0, 140, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 200, 47,
              ]

  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);
  # foo.extend(foo);

  bar = []
  bar.extend(pre)
  bar.extend(foo)
  # bar.extend(pre)
  bar.extend([0, 140, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 94, 0, 47,
              0, 94, 0, 47,
              0, 47, 0, 47,
              0, 47, 0, 47,
              0, 47, 8, 47,
              ])

  t.transmit(bar)

