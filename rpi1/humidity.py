#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from sense_hat import SenseHat
from time import sleep
from api_manager import update_values_api
import os, math, time

def get_hum():
    # Check comments on file thermometre.py, same logic
    sense = SenseHat()
    hum = round(sense.get_humidity(),1)
    last_hum = -100.0
    if os.path.isfile("/tmp/last_hum.txt"):
        last_hum_f = open("/tmp/last_hum.txt", "r")
        last_hum = float(last_hum_f.readline())
        last_hum_f.close()
    tries = 0

    while abs(last_hum - hum) > 4 and tries < 3:
        hum = round(sense.get_humidity(), 1)
        hum = round(sense.get_humidity(), 1)
        if abs(last_hum - hum) > 4:
            sleep(5)
            tries += 1

    if last_hum < -10:
        os.system("echo "+str(hum)+" > /tmp/last_hum.txt")
        return hum
    if math.isnan(hum) or hum <= 0.0 or abs(last_hum - hum) > 4 or hum > 100.0:
        os.system("echo "+str(last_hum)+" > /tmp/last_hum.txt")
        update_values_api({"Hum":last_hum})
        return last_hum
    os.system("echo "+str(hum)+" > /tmp/last_hum.txt")
    update_values_api({"Hum":hum})
    return hum


if __name__ == "__main__":
    while True:
        get_hum()
        time.sleep(60)

