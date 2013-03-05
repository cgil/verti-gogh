'''
Created on Mar 4, 2013

@author: carlosgil
'''

from Classes import Car
from Classes import VertCanvas
from Tkinter import Tk, BOTH
from random import randint
        
root = Tk()
width = 800
height = 600
root.geometry(str(width)+"x"+str(height))
canvas = VertCanvas.VertCanvas(root)
cars = []
for i in range(2) : 
    cars.append(Car.Car(i,randint(1,width),randint(1,height), canvas))


def draw():
    for car in cars :
        car.addPoint(randint(1,width), randint(1,height))
        car.display()
    canvas.pack(fill=BOTH, expand=1)
    root.after(1000, draw)

root.after(1000, draw)
root.mainloop()  


