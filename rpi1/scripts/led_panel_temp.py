#!/usr/bin/env python
# -*- coding: utf-8 -*-

from sense_hat import SenseHat
from time import sleep
from threading import Thread

from controllers.display import update_display_matrix, colour
from thermometre import get_temp


def update_temp(temp, ledDisplay):
    last_temp = 0
    while True:
        temp = int(round(get_temp()))
        if temp != last_temp:
            update_display_matrix(temp, ledDisplay)
            sense.set_pixels(ledDisplay)
        
        if temp >= 30:
            blink(&temp, None, None)

        last_temp = temp
        sleep(10)


def blink(temp, ledDisplay, ledDisplayOff):
    while True:
        while(temp > 30):
            sense.set_pixels(ledDisplay)
            sleep(0.4)
            sense.set_pixels(ledDisplayOff)
            sleep(0.2)


def main():
    # Initialize SenseHat
    sense = SenseHat()
    sense.set_rotation(270) # Rotate the led display axis to fit its position
    sense.low_light = True # Low brightness mode because the room sometimes is darker

    # Initialize led panel
    _, O = colour(0) # Get the initial RGB colour config (all pixels off)
    ledDisplay = [] # Image when led panel is on
    ledDisplayOff = [] # Image when led panel is off (all pixels off)
    for _ in range(64):
        ledDisplay.append(O)
        ledDisplayOff.append(O)


    t_update_temp = Thread(target=update_temp,args=(ledDisplay))
    t_update_temp.start()

    t_blink = Thread(target=blink,args=(temp,ledDisplay,ledDisplayOff))
    t_blink.start()

    # Wait for thread
    t_update_temp.join()
    t_update_temp.join()


if __name__ == "__main__":
    main()