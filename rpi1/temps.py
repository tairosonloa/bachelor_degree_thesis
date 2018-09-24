#!/usr/bin/env python

from sense_hat import SenseHat
from time import sleep
from threading import Thread

from definition.colours import colour
from thermometre import get_temp
from definition.display import display


# Initialize SenseHat
sense = SenseHat()
sense.set_rotation(270)
sense.low_light = True

X, O = colour(0)

ledDisplay = []
ledDisplayOff = []
for i in range(64):
    ledDisplay.append(O)
    ledDisplayOff.append(O)
temp = [0]
sense.set_pixels(ledDisplay)


def updateTemp(temp, ledDisplay):
    aux = 0
    while(1):
        temp[0] = int(round(get_temp()))
        if(temp[0] != aux):
            display(temp[0], ledDisplay)
            sense.set_pixels(ledDisplay)
        aux = temp[0]
        sleep(10)

def parpadeo(temp, ledDisplay, ledDisplayOff):
    while(1):
        while(temp[0] > 30):
            sense.set_pixels(ledDisplay)
            sleep(0.4)
            sense.set_pixels(ledDisplayOff)
            sleep(0.2)

        sleep(10)
thread = Thread(target=updateTemp,args=(temp,ledDisplay))
thread.start()

thread2 = Thread(target=parpadeo,args=(temp,ledDisplay,ledDisplayOff))
thread2.start()
thread.join()
thread2.join()

