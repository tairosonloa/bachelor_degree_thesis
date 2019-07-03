# install/
This folder contains all files used by the installation script to deploy the system.

## rpi1/daemons/
Systemd daemons used to launch the rpi1 python scripts on system bootup.

## rpi2/
Daemons, scripts and binaries used to deploy rpi2. More information is inside that folder.

## rpi3/
Daemons, scripts and binaries used to deploy rpi3. More information is inside that folder.

## .bashrc
Bash configuration file used on the Raspberry Pi computers. It coaints a colored prompt, ls colors, and ls aliases

## default.config
Example of file with default configuration values, used by the installation script to automatize or speed up the installation process.

## .raspi-monitor
A bash script used by the installation script to automatically switch on/off the monitor of rpi2 and rpi3 accordingly with the laboratory working hours.