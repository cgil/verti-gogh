'''
Created on Mar 4, 2013

@author: carlosgil
'''

import sys
import time

from Classes import Car
from Classes import VertCanvas
from Tkinter import Tk, BOTH, tkinter
from random import randint
from subprocess import Popen, PIPE, STDOUT

root = Tk()
width = 1024
height = 768
root.geometry(str(width)+"x"+str(height))
canvas = VertCanvas.VertCanvas(root)
cars = []
for i in range(1) :
  cars.append(Car.Car(i,randint(1,width),randint(1,height), canvas))

topx, topy, botx, boty = 0, 0, 0, 0

def clear():
  canvas.canvas.create_polygon(0, 0, width, 0, width, height, 0, height,
                               fill='white')

def circle(x, y, rad, color):
  canvas.canvas.create_oval(x - rad, y - rad, x + rad, y + rad,
                            fill=color, outline=color)

colors = []
buckets = []
color = None
proc = None

def calibrate_start():
  global proc
  proc = Popen(['./capture/ping'], stdout=PIPE, stdin=PIPE, stderr=None)
  calibrate()

def calibrate():
  global colors, buckets, color
  print colors, buckets
  for i, clrs in enumerate(colors):
    if len(clrs) >= 5:
      print "found at ", buckets[i][0], buckets[i][1]
      proc.terminate()
      proc.wait()
      return

  color = '%06x' % randint(0, (1 << 24) - 1)
  clear()
  circle(width / 2, height / 2, 50, '#' + color)
  canvas.pack(fill=BOTH, expand=1)
  root.after(400, check)

def check():
  global colors, buckets, proc
  proc.poll()
  if proc.returncode != None and proc.returncode != 0:
    print("bad return code", proc.returncode)
    raise 'oh no'

  proc.stdin.write('0x' + color + "\n")
  x, y = proc.stdout.readline()[:-1].split(' ')

  appended = False
  for i, bkt in enumerate(buckets):
    dx = bkt[0] - int(x)
    dy = bkt[1] - int(y)
    d = dx * dx + dy * dy
    if d < 4:
      colors[i].append(color)
      appended = True
      break
  if not appended:
    colors.append([color])
    buckets.append([int(x), int(y)])

  calibrate()



def readappend(fh, _):
  mystr = sys.stdin.readline()
  print mystr
  idx, row, col = mystr.split(' ')
  cars[int(idx)].addPoint(int(row), int(col))
  cars[int(idx)].display()
  canvas.pack(fill=BOTH, expand=1)

root.after(10, calibrate_start)

root.mainloop()
