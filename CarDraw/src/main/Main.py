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
height = 500
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

curx = None
cury = None

center = None
topleft = None
topright = None
botleft = None
botright = None
painted = False

# Check the buckets if one of the contains the "true location of the dot".
def check():
  global colors, buckets, color
  print colors, buckets
  for i, clrs in enumerate(colors):
    if len(clrs) >= 5:
      return i
  return None

# Locate the dot, put it in a bucket
def locate():
  global colors, buckets, color
  proc = Popen(['./capture/find_raw', '0x' + color],
               stdout=PIPE, stdin=None, stderr=None)

  x, y = proc.stdout.readline()[:-1].split(' ')
  proc.wait()
  if proc.returncode != None and proc.returncode != 0:
    print("bad return code", proc.returncode)
    raise 'oh no'
  print color, 'at', x, y

  appended = False
  for i, bkt in enumerate(buckets):
    dx = bkt[0] - int(x)
    dy = bkt[1] - int(y)
    d = dx * dx + dy * dy
    if d < 36:
      colors[i].append(color)
      appended = True
      break
  if not appended:
    colors.append([color])
    buckets.append([int(x), int(y)])

def calibrate_center():
  global color, center, painted
  if painted:
    locate()
  i = check()
  if i == None:
    painted = True
    color = '%06x' % randint(0, (1 << 24) - 1)
    clear()
    circle(width / 2, height / 2, 20, '#' + color)
    canvas.pack(fill=BOTH, expand=1)
    root.after(200, calibrate_center)
  else:
    center = buckets[i]
    print 'found center', center
    color = colors[i][0]
    painted = False
    calibrate_corners()

def calibrate_corners():
  global topleft, topright, botleft, botright, color, painted, buckets, colors
  if painted:
    locate()
    if topleft == None:
      topleft = buckets[0]
    elif topright == None:
      topright = buckets[0]
    elif botleft == None:
      botleft = buckets[0]
    else:
      botright = buckets[0]
      done_calibration()
      return
  buckets = []
  colors = []
  painted = True
  clear()
  if topleft == None:
    circle(20, 20, 20, '#' + color)
  elif topright == None:
    circle(width - 20, 20, 20, '#' + color)
  elif botleft == None:
    circle(20, height - 20, 20, '#' + color)
  else:
    circle(width - 20, height - 20, 20, '#' + color)
  canvas.pack(fill=BOTH, expand=1)
  root.after(200, calibrate_corners)

mystdin = None

bounds = []

def done_calibration():
  global mystdin, colors, buckets, center, topleft, topright, botleft, botright, bounds
  print 'center', center
  print 'topleft', topleft
  print 'topright', topright
  print 'botleft', botleft
  print 'botright', botright
  bad = False
  if topleft[0] > center[0] or topleft[1] > center[1]:
    print 'invalid topleft'
    bad = True
  if topright[0] < center[0] or topright[1] > center[1]:
    print 'invalid topright'
    bad = True
  if botleft[0] > center[0] or botleft[1] < center[1]:
    print 'invalid botleft'
    bad = True
  if botright[0] < center[0] or botright[1] < center[1]:
    print 'invalid botright'
    bad = True

  if bad:
    print 'make sure the webcam sees the whole screen, retrying'
    painted = False
    colors = []
    buckets = []
    topleft = topright = botleft = botright = center = None
    root.after(1000, calibrate_center)
    return
  bounds = [[min(topleft[0], botleft[0]),
             min(topleft[1], topright[1])],
            [max(botright[0], topright[0]),
             max(botright[1], botleft[1])]]

  proc = Popen(['./capture/capture_raw_frames',
                '0x000000',
                str(bounds[0][0]),
                str(bounds[0][1]),
                str(bounds[1][0]),
                str(bounds[1][1])],
               stdout=PIPE)
  mystdin = proc.stdout
  root.tk.createfilehandler(mystdin, tkinter.READABLE, readappend)

prevrow = 0
prevcol = 0
def readappend(fh, _):
  global prevrow, prevcol
  mystr = mystdin.readline()
  # print mystr
  row, col = mystr.split(' ')
  row = (int(row) - bounds[0][1]) * height / (bounds[1][1] - bounds[0][1])
  col = (int(col) - bounds[0][0]) * width / (bounds[1][0] - bounds[0][0])

  circle(prevcol, prevrow, 5, '#ffffff')
  circle(col, row, 5, '#ff0000')
  canvas.pack(fill=BOTH, expand=1)
  prevrow = row
  prevcol = col

clear()
root.after(10, calibrate_center)
root.overrideredirect(True)
root.mainloop()
