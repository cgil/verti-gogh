'''
Created on Mar 4, 2013

@author: carlosgil
'''

import sys

from Classes import Car
from Classes import VertCanvas
from Tkinter import Tk, BOTH, tkinter
from random import randint
import subprocess

root = Tk()
width = 1024
height = 768
root.geometry(str(width)+"x"+str(height))
canvas = VertCanvas.VertCanvas(root)
cars = []
for i in range(1) :
    cars.append(Car.Car(i,randint(1,width),randint(1,height), canvas))

topx, topy, botx, boty = 0, 0, 0, 0

def draw():
    for car in cars :
        car.addPoint(randint(1,width), randint(1,height))
        car.display()
    canvas.pack(fill=BOTH, expand=1)
    root.after(1000, draw)

def circle(x, y, rad):
  canvas.canvas.create_oval(x - rad, y - rad, x + rad, y + rad, fill='black')

def drawcorners():
    rad = 10
    circle(rad, rad, rad)
    circle(width - rad, rad, rad)
    circle(rad, height - rad, rad)
    circle(width - rad, height - rad, rad)
    canvas.pack(fill=BOTH, expand=1)
    root.after(1000, calibrate)

def trans(x, y):
    return (width * (x - topx) / (botx - topx),
            height * (y - topy) / (boty - topy))

def calibrate():
    global topx, topy, botx, boty
    output = subprocess.check_output(["./capture/calibrate_raw_frame"])
    c1, c2, c3, c4, _ = output.split("\n")
    c1x, c1y = c1.split(' ')
    c2x, c2y = c2.split(' ')
    c3x, c3y = c3.split(' ')
    c4x, c4y = c4.split(' ')

    topx = (int(c1x) + int(c3x)) / 2
    topy = (int(c1y) + int(c2y)) / 2
    botx = (int(c2x) + int(c4x)) / 2
    boty = (int(c3y) + int(c4y)) / 2

    # TODO: run capture raw frames and hook it up to readappend

def readappend(fh, _):
    mystr = sys.stdin.readline()
    print mystr
    idx, row, col = mystr.split(' ')
    cars[int(idx)].addPoint(int(row), int(col))
    cars[int(idx)].display()
    canvas.pack(fill=BOTH, expand=1)

if len(sys.argv) > 1 and sys.argv[1] == 'rand':
  root.after(1000, draw)
# else:
# root.tk.createfilehandler(sys.stdin, tkinter.READABLE,
#                           readappend)
root.after(10, drawcorners)

root.mainloop()
