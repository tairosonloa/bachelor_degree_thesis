#! /usr/bin/python3
# -*- coding: utf-8 -*-

from controllers.api_manager import update_values_api

import numpy as np
import picamera, picamera.array
import time, os, math, json


SECONDS_BETWEEN_MEASURES = 10 # Time between two measurements
TRESHOLD = 50 # Brightness treshold to distinguish light off from light on


def is_cpd_light_on():
    """Checks if light inside CPD room is on"""
    # Initialize Pi camera
    with picamera.PiCamera() as camera:
        camera.resolution = (128,80)
        with picamera.array.PiRGBArray(camera) as stream:
            # Configure camera
            camera.exposure_mode = 'auto'
            camera.awb_mode = 'auto'
            # Take a photo (of the wall)
            camera.capture(stream, format='rgb')
            # Check the pixels brightness average
            pixAverage = int(np.average(stream.array[...,1]))
    # Based on pixels brightness average, we indicate is light is on or off
    if(pixAverage > TRESHOLD):
        return True # ON
    else:
        return False # OFF


def main():
    """Checks light status every 10 seconds and updates API"""
    # Check light status every 10s and update API
    while True:
        if is_cpd_light_on():
            update_values_api({ "Light" : True })
        else:
            update_values_api({ "Light" : False })
        time.sleep(SECONDS_BETWEEN_MEASURES)


if __name__ == "__main__":
    main()
