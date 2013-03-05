'''
Created on Mar 4, 2013

@author: carlosgil
'''

import sys

from Classes import Car
from Classes import VertCanvas
from Tkinter import Tk, BOTH, tkinter
from random import randint
        
root = Tk()
width = 320
height = 240
root.geometry(str(width)+"x"+str(height))
canvas = VertCanvas.VertCanvas(root)
cars = []
for i in range(1) :
    cars.append(Car.Car(i,randint(1,width),randint(1,height), canvas))


def draw():
    for car in cars :
        car.addPoint(randint(1,width), randint(1,height))
        car.display()
    canvas.pack(fill=BOTH, expand=1)
    root.after(1000, draw)

def readappend(fh, _):
    mystr = sys.stdin.readline()
    print mystr
    idx, row, col = mystr.split(' ')
    cars[int(idx)].addPoint(int(row), int(col))
    cars[int(idx)].display()
    canvas.pack(fill=BOTH, expand=1)

if len(sys.argv) > 1 && sys.argv[1] == 'rand'
  root.after(1000, draw)
else
  root.tk.createfilehandler(sys.stdin, tkinter.READABLE,
                            readappend)

root.mainloop()  
