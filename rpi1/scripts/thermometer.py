#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from controllers.CPUTemp import CPUTemp
from controllers.api_manager import update_values_api

from sense_hat import SenseHat
from time import sleep
import os, math


FILE_PATH = "/tmp/last_temp.txt" # Where last measurement is/will be stored
# If readed temp differs from last temp in more than 1 celsius degree we treat it
# as a measurement error, so we will wait 5 seconds and try again for 3 times
SECONDS_BETWEEN_MEASURES = 60
ERROR_TRESHOLD = 1
TRIES = 3
WAIT_SECONDS = 5


def update_state(value):
    """Saves readed measurement as last measurement and updates API"""
    # Save on file last measurement to check measuring errors
    with open(FILE_PATH, "w") as f:
        f.write(str(value))
    # POST request to API on rpi2 to update values
    update_values_api({"Temp":value})


def get_temp():
    """Reads current temperature from sensor"""
    # Initialize SenseHat
    sense = SenseHat()
    temp = 0.0
    
    # Check if last measurement exists, used to detect possible measurement error when call the sensor
    last_temp = -100.0
    if os.path.isfile(FILE_PATH):
        last_temp_f = open(FILE_PATH, "r")
        last_temp = float(last_temp_f.readline())
        last_temp_f.close()
    tries = 0 # If detected measurement error, tries counter to accept the measurement as valid
    
    # Algorithm to adjust impact of the CPU temp on the temperature sensor readings
    while abs(last_temp - temp) > ERROR_TRESHOLD and tries < TRIES:
        # We need two continuous measures because first measurement uses to fail
        p = sense.get_temperature_from_pressure()
        h = sense.get_temperature_from_humidity()
        p = sense.get_temperature_from_pressure()
        h = sense.get_temperature_from_humidity()
        
        # Calculates temperature
        with CPUTemp() as cpu_temp:
            c = cpu_temp.get_temperature()
        temp = round(((p+h)/2) - (c/7),1)

        # Check if possible measurement error and wait 5 seconds to try again
        if abs(last_temp - temp) > ERROR_TRESHOLD:
            sleep(WAIT_SECONDS)
            tries += 1

    # If we did a wrong measurement, use the last measurement as current measurement
    # (if last_temp == -100, first measure from reboot)
    if (math.isnan(temp) or temp == 0.0 or abs(last_temp - temp) > ERROR_TRESHOLD) and last_temp != -100:
        update_state(last_temp) # Save measurement and call API to update
        return last_temp
    update_state(temp) # Save measurement and call API to update
    return temp


# If called as standalone, check humidity every 60 seconds
if __name__ == "__main__":
    while True:
        get_temp()
        sleep(SECONDS_BETWEEN_MEASURES)
