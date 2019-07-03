#!/usr/bin/env python
# -*- coding: utf-8 -*-

from thermometer import get_temp
from controllers.Display import Display


from time import sleep

SECONDS_BETWEEN_MEASURES = 30

def main():
    """Show temperature in display and refresh every 10 seconds"""
    display = Display()
    while True:
        display.update_display(round(get_temp()))
        sleep(SECONDS_BETWEEN_MEASURES)

if __name__ == "__main__":
    main()