'''
Created on Mar 4, 2013

@author: carlosgil
'''
from Tkinter import Frame, BOTH, Canvas

class VertCanvas(Frame):
  
    def __init__(self, parent):
        Frame.__init__(self, parent, background="white")   
        self.parent = parent
        self.canvas = Canvas(self)         
        self.canvas.pack(fill=BOTH, expand=1)
