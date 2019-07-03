# rpi1/scripts/controllers/
This folder contains some controllers and auxiliary functions used by the scripts.

## api_manager.py
Makes an authenticated HTTP POST request to `rpi2_api` to update some value. It uses `config_loader.py` to load the configuration needed for the API requests.

## config_loader.py
Loads the JSON configuration file from `etc/rpi1_conf.json`

## CPUTemp.py
Python object which contains all methods needed to avoid Raspbery Pi CPU heat interfering on the temperature measures taken for the data center.

## Display.py
Python object which contains all methods needed to control the SenseHat led panel.