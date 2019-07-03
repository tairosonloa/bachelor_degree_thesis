# install/rpi1/daemons/
This folder contains the systemd daemons used to launch the rpi1 python scripts on system bootup.

## cpd_humidity.service
Automatically launchs `hygrometer.py` on system bootup.

## cpd_light.service
Automatically launchs `light_watcher.py` on system bootup.

## cpd_temperature.service
Automatically launchs `led_display_temp.py` on system bootup.