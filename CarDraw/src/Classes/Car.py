'''
Created on Mar 4, 2013

@author: carlosgil
'''

from Classes import Point
from random import randint

class Car(object):

    def __init__(self, _index, _x, _y, _vertCanvas):
        self.index = _index
        self.points = []
        self.getColor()
        self.addPoint(_x, _y)
        self.vertCanvas = _vertCanvas
        
    def addPoint(self, _x, _y):
        newPoint = Point.Point(_x, _y)
        self.points.append(newPoint)
        
    def display(self):
        if self.points.__len__() > 0 :
            radius = 10
            lastPoint = self.points[0]
            self.vertCanvas.canvas.create_oval(lastPoint.x-radius, lastPoint.y-radius, lastPoint.x+radius, lastPoint.y+radius, width=2, fill='blue')
            if self.points.__len__() > 1 :
                for p in self.points :
                    self.vertCanvas.canvas.create_oval(p.x-radius, p.y-radius, p.x+radius, p.y+radius, width=2, fill=self.color)
                    self.vertCanvas.canvas.create_line(lastPoint.x, lastPoint.y, p.x, p.y, fill=self.color, dash=(4, 4))
                    lastPoint = p
        
    def getColor(self):
        r = randint(1, 254)
        g = randint(1, 254)
        b = randint(1, 254)
        rgb = r, g, b
        self.color = "#%02x%02x%02x" % rgb

        