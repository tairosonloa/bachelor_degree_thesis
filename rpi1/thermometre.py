#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from sense_hat import SenseHat
from definition.CPUTemp import CPUTemp
from time import sleep
from api_manager import update_values_api
import os, math

def get_temp():
    # Initialize SenseHat
    sense = SenseHat()

    temp = 0.0
    last_temp = -100.0
    # Rpi was shut down, so last measure isn't trustworthy
    if os.path.isfile("/tmp/last_temp.txt"):
        last_temp_f = open("/tmp/last_temp.txt", "r")
        last_temp = float(last_temp_f.readline())
        last_temp_f.close()
    tries = 0 
    
    # Algorithm to adjust impact of the CPU temp on the temperature sensor readings
    while abs(last_temp - temp) > 1 and tries < 3:
        p = sense.get_temperature_from_pressure()
        h = sense.get_temperature_from_humidity()
        p = sense.get_temperature_from_pressure()
        h = sense.get_temperature_from_humidity()
        with CPUTemp() as cpu_temp:
            c = cpu_temp.get_temperature()
        temp = round(((p+h)/2) - (c/7),1)
        if abs(last_temp - temp) > 1:
            sleep(5)
        tries += 1

    # Write a temperature less than -10 on file last_temp.txt to manually restart
    if last_temp < -10:
        os.system("echo "+str(temp)+" > /tmp/last_temp.txt")
        return temp
    if math.isnan(temp) or temp == 0.0 or abs(last_temp - temp) > 1:
        # Save on file last measurement to check measuring errors
        os.system("echo "+str(last_temp)+" > /tmp/last_temp.txt")
        # POST request to API on rpi2 to update values
        update_values_api({"Temp":last_temp})
        return last_temp
    # Save on file last measurement to check measuring errors
    os.system("echo "+str(temp)+" > /tmp/last_temp.txt")
    # POST request to API on rpi2 to update values
    update_values_api({"Temp":temp})
    return temp

if __name__ == "__main__":
    print(get_temp()) 

