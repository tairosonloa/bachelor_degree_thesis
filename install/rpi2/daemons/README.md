# install/rpi2/daemons/
This folder contains the systemd daemons used to launch the rpi2 API REST on system bootup.

## rpi2_api.service
Automatically launchs `rpi2_api_arm` binary on system bootup and sets the binary to use `/etc/rpi2_conf.json` as configuration file.