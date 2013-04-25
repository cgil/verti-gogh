#!/usr/bin/env python
#
# Class for simplifying the reading and transmitting of IR signals from the IR Toy.
# This only works for firmware revision 22 or greater.
# see https://github.com/crleblanc/PyIrToy and
# http://dangerousprototypes.com/docs/USB_Infrared_Toy for more info.
#
# Chris LeBlanc, 2012
#
#--
#
# This work is free: you can redistribute it and/or modify it under the terms
# of Creative Commons Attribution ShareAlike license v3.0
#
# This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
# without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
# See the License for more details. You should have received a copy of the License along
# with this program. If not, see <http://creativecommons.org/licenses/by-sa/3.0/>.

import sys
import time
import binascii
import logging

__author__ = 'Chris LeBlanc'
__version__ = '0.2.1'
__email__ = 'crleblanc@gmail.com'

class IrToy(object):

    def __init__(self, serialDevice):
        self.toy = serialDevice
        logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
        # ir toy must be in sampling mode, even for transmissions to work
        self._setSamplingMode()

    def _setSamplingMode(self):
        '''set the IR Toy to use sampling mode, which we use exclusively'''
        self.reset()
        logging.debug('reset')
        self.write(b'S')
        logging.debug('wrote sampling mode')
        protocolVersion = self.toy.read(3)
        if protocolVersion != 'S01':
          raise IOError('Expected protocol "S01", got %s' % protocolVersion)
        logging.debug('got protocol %s' % protocolVersion)

    def _acknowledge(self, code):
        byteCode = bytearray(code)
        for idx in range(0, len(code), 62):
          logging.debug('sending %r', code[idx:idx+62])
          self.write(byteCode[idx:idx+62])
          ack = ord(self.toy.read(1))
          if ack != 0x3e:
            raise IOError('invalid acknowledgement received %d' % ack)


    def reset(self):
      logging.debug('resetting')
      self.write(bytearray([0x00] * 5))

    def write(self, arr):
      # time.sleep(0.05)
      self.toy.write(arr)

    def transmit(self, code):
        if len(code) % 2 != 0:
            raise ValueError("Length of code argument must be an even number")

        # ensure the last two codes are always 0xff (255) to tell the IR Toy
        # it's the end of the signal
        if code[-2:] != [0xff, 0xff]:
            code.extend([0xff, 0xff])

        # Enable transmit handshake, notify on complete, byte count report, and
        # then sent the actual transmit command
        self._acknowledge([0x26, 0x25, 0x24, 0x03])
        logging.debug('sent transmit')

        # Send the actual data (chunked inside)
        self._acknowledge(code)
        logging.debug('sent code')

        # Acknowledge the data was actually sent, and all of it was sent.
        if self.toy.read(1) != 't':
          raise IOError('expected a "t" acknowledgement')
        hexBytes = binascii.b2a_hex(self.toy.read(2))
        if self.toy.read(1) != 'C':
          raise IOError('expected a completion "C"')
        if int(hexBytes, 16) != len(code):
          raise IOError('short write')

        # self._setSamplingMode()
