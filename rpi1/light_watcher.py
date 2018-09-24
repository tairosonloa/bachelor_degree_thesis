#! /usr/bin/python3

from phue import Bridge
from api_manager import update_values_api
import time, os, math
import picamera
import picamera.array
import numpy as np
import json

def is_cpd_light_on():
    """Comprueba si la luz del CPD está encendida en un momento dado"""
    with picamera.PiCamera() as camera:
        camera.resolution = (128,80)
        with picamera.array.PiRGBArray(camera) as stream:
            camera.exposure_mode = 'auto'
            camera.awb_mode = 'auto'
            camera.capture(stream, format='rgb')
            pixAverage = int(np.average(stream.array[...,1]))
    if(pixAverage > 50):
        return True
    else:
        return False


def cpd_light_watcher():
    """Comprueba cada 10 si la luz del CPD está encendida"""
    while True:
        if is_cpd_light_on():
            update_values_api({ "Light" : True })
        else:
            update_values_api({ "Light" : False })
        time.sleep(10)

config = {}
with open("config.json", "r") as f:
        config = json.load(f)
b = Bridge(config["HueBridgeAddress"])
b.connect()
cpd_light_watcher()
