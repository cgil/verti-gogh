import serial
import time
import sys

debug = False
toy = None

def write(input):
    time.sleep(0.05)
    tmp = toy.write(chr(input))
    if (debug):
        print "\tBytes written: %d" % tmp

def reset():
    print "[*] Resetting "
    write(0x00)
    write(0x00)
    write(0x00)
    write(0x00)
    write(0x00)
    print "[*] Reset done"

    
if __name__ == "__main__":
    toy = serial.Serial("/dev/ttyACM0", 9600)
        
    reset()
    toy.write('S')    #enter sampling mode
    time.sleep(0.05)
    print "[*] Protocol version %s" % toy.read(3)

    write(0x26)             #Enable transmit handshake
    write(0x25)             #Enable transmit notify on complete
    write(0x24)             #Enable transmit byte count report
    
    write(0x03)             #Start transmit mode
    print "[*] Handshake: %s" % hex(ord(toy.read(1)))
    
    x = bytearray([0x11,0x11,0xff,0xff])
    toy.write(x)
    print "[*] Final handshake: %s" % hex(ord(toy.read(1)))
   
    print "[*] Transmit byte count: %s" % (toy.read(3)[1:].encode("HEX"))
    print "[*] Notify on complete: %s" % toy.read(1)
    reset()
    
    print "[*] Closing..."
    toy.close()
    print "[*] Closed"
    print "[*] Trying to open again"
    toy = serial.Serial("/dev/ttyACM0", 9600)
    print "[*] Success!"

