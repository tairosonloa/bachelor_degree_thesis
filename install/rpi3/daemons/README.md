# install/rpi3/daemons/
This folder contains the systemd daemons used to launch the rpi3 API REST and custom light web server on system bootup.

## rpi3_api.service
Automatically launchs `rpi3_api_arm` binary on system bootup and sets the binary to use `/etc/rpi3_conf.json` as configuration file.

## rpi3_web_server.service
Automatically launchs `web_server_arm` binary on system bootup with default configuration.