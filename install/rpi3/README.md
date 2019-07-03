# install/rpi3/
This folder contains the daemons, scripts and binaries used to deploy rpi3.

## autostart
Bash script used to automatically start chromium browser with rpi3 GUI running on localhost on full screen.

## rpi3_api_arm
Rpi3 API REST binary compiled for ARM processor architecture. Can be called with the following parameters:
* `-conf`: path tho the JSON configuration file (default `$(pwd)/config.json`).

## web_Server_arm
Custom light web server binary compiled for ARM processor architecture. Can be called with the following parameters:
* `-port`: port where web server will be listening (default 9000).
* `-root`: rpi3 GUI static files root path (default `/srv/rpi3/`).