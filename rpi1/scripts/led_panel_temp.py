#!/usr/bin/env python
# -*- coding: utf-8 -*-

from thermometre import get_temp
from controllers.Display import Display


from time import sleep


def main():
    """Show temperature in display and refresh every 10 seconds"""
    display = Display()
    while True:
        display.update_display(round(get_temp()))
        sleep(10)

if __name__ == "__main__":
    main()