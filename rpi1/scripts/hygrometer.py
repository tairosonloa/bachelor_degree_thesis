#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from controllers.api_manager import update_values_api

from sense_hat import SenseHat
from time import sleep
import os, math


FILE_PATH = "/tmp/last_hum.txt" # Where last measurement is/will be stored
# If readed hum differs from last hum in more than 4 % we treat it as a
# measurement error, so we will wait 5 seconds and try again for 3 times
SECONDS_BETWEEN_MEASURES = 60
ERROR_TRESHOLD = 4
TRIES = 3
WAIT_SECONDS = 5


def update_state(value):
    """Saves readed measurement as last measurement and updates API"""
    # Save on file last measurement to check measuring errors
    with open(FILE_PATH, "w") as f:
        f.write(str(value))
    # POST request to API on rpi2 to update values
    update_values_api({"Hum":value})


def get_hum():
    """Reads current humidity from sensor"""
    # Initialize SenseHat
    sense = SenseHat()
    hum = 0.0

    # Check if last measurement exists, used to detect possible measurement error when call the sensor
    last_hum = -100.0
    if os.path.isfile(FILE_PATH):
        last_hum_f = open(FILE_PATH, "r")
        last_hum = float(last_hum_f.readline())
        last_hum_f.close()
    tries = 0 # If detected measurement error, tries counter to accept the measurement as valid

    # Read the current humidity from sensor
    while abs(last_hum - hum) > ERROR_TRESHOLD and tries < TRIES:
        # We need two continuous measures because first measurement uses to fail
        hum = round(sense.get_humidity(), 1)
        hum = round(sense.get_humidity(), 1)
        # Check if possible measurement error and wait 5 seconds to try again
        if abs(last_hum - hum) > ERROR_TRESHOLD:
            sleep(WAIT_SECONDS)
            tries += 1

    # If we did a wrong measurement, use the last measurement as current measurement
    # (if last_hum == -100, first measure from reboot)
    if (math.isnan(hum) or hum <= 0.0 or abs(last_hum - hum) > 4 or hum > 100.0) and last_hum != -100:
        update_state(last_hum) # Save measurement and call API to update
        return last_hum
    update_state(hum) # Save measurement and call API to update
    return hum


# If called as standalone, check humidity every 60 seconds
if __name__ == "__main__":
    while True:
        get_hum()
        sleep(SECONDS_BETWEEN_MEASURES)
