# rpi1/scripts/
This folders contains all the python 3 scripts used to monitor the data center.

## controllers/
Contains some controllers and auxiliary functions used by the scripts.

## models/
Contains custom data structures models used in the scripts.

## hygrometer.py
Checks humidity inside the data center every 60 seconds.

## led_display_temp.py
Checks temperature inside the data center every 30 seconds (by calling `thermometer.py`) and updates the value on the SenseHat led panel accordingly.

## light_watcher.py
Checks light status (on/off) inside the data center every 10 seconds.

## thermometer.py
Checks temperature inside the data center every 60 seconds if runned standalone. Used by `led_display_temp.py` to check the temperature every 30 seconds.