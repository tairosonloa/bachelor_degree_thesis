#! /usr/bin/env python3

import sys, os


def create_config_json():
    pass


def install_rpi1():
    print("##### Installing dependences with apt-get...")
    # Install python libs
    print("##### Enabling components on rpi...")
    # Enable camera and ic2
    print("##### Preparing daemons and start on boot...")
    # Copy daemons and enable it
    print("##### Setting and enabling iptables...")
    # Set and enable iptables
    print("##### Generating config.json...")
    create_config_json()
    print("##### Starting scripts...")


def install_rpi2():
    print("##### Copying files...")
    # Copy binary and alarm.mp3
    print("##### Preparing daemons and start on boot...")
    # Copy daemons and enable it
    print("##### Setting and enabling iptables...")
    # Set and enable iptables
    print("##### Generating config.json...")
    create_config_json()
    input("##### WARNING: Please, press the philips hue bridge button before continue. Press enter when pressed.")
    print("##### Starting scripts...")


print("What rpi are you trying to install?\n\n" +
        "    (1) rpi1\n" +
        "    (2) rpi2\n")
opt = int(input("Please choose a number. Choose (0) to exit: "))
print()

if opt == 0:
    sys.exit(0)
elif opt == 1:
    install_rpi1()
elif opt == 2:
    install_rpi2()

print("\nDone. If you see this message, everything should be working now.")