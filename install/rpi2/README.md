# install/rpi2/
This folder contains the daemons, scripts and binaries used to deploy rpi2.

## alarm.mp3
Siren MP3 sound used to fire audio alerts. [CC-By 3.0 Mike Koening](http://soundbible.com/287-Industrial-Alarm.html).

## autostart
Bash script used to automatically start chromium browser with grafana instance running on ultraheroe on full screen.

## rpi2_api_arm
Rpi2 API REST binary compiled for ARM processor architecture. Can be called with the following parameters:
* `-conf`: path to the JSON configuration file (default `$(pwd)/config.json`).