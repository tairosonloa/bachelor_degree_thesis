# -*- coding: utf-8 -*-

from models.digits import digits

from sense_hat import SenseHat
from threading import Thread
from time import sleep


PIXEL_OFF = [0,0,0] # RGB for pixel off


def colour_by_temp(temp):
    """Returns rgb color code list based on temperature range and pixel off rgb list"""
    
    if(temp >= 40):
        X = [255, 0, 0]     # Red
    elif(temp >= 30):
        X = [255, 128, 0]   # Orange
    elif(temp >= 20):
        X = [255, 255, 0]   # Yellow
    elif(temp >= 10):
        X = [0, 255, 0]     # Green
    elif(temp >= 0):
        X = [0, 255, 128]   # Turquoise
    elif(temp >= -10):
        X = [0, 255, 255]   # Light blue
    elif(temp >= -20):
        X = [0, 191, 255]   # Blue
    elif(temp >= -30):
        X = [0, 127, 255]   # Darker blue
    elif(temp >= -40):
        X = [0, 64, 255]    # Dark blue
    else:
        X = [0, 0, 255]     # Very dark blue

    O = PIXEL_OFF # Black == off state

    return X,O


class Display:
    def __init__(self):
        self.sense = SenseHat()
        self.sense.set_rotation(270) # Rotate the led display axis to fit its position
        self.sense.low_light = True # Low brightness mode because the room sometimes is darker
        
        self.temp = 0 # Temperature currently beeing displaying
        self.blinking = False # DIsplay is blinking or not

        # Current display and display off
        self.pixels_matrix = [(0,0,0)] * 64     # Pixels matrix currently being displaying
        self.pixels_matrix_off = ((0,0,0)) * 64 # Pixels matrix with all pixels off
    
    def update_display(self, temp):
        """Refresh the pixels matrix display based on temperature to be displayed"""
        # If same temperature, it's unnecesary to refresh the display
        if self.temp != temp:
            # Update last temperature
            self.temp = temp
            # Update colors to fit the new temperature
            X, O = colour_by_temp(temp)

            # Since the display is only big enough for two digits, an exception is made for anything over 99
            # If the temp does hit +/- 100 then it will blank out the display given that it still works
            if abs(temp) >= 100:
                for i in range(64):
                    self.pixels_matrix[i] = X
            else:
                # Start building the display array (pixels_matrix)
                index = 0
                digitIndex = 0
                digits_representation = digits(X, O)
                left_digit = int(abs(temp) / 10)
                right_digit = int(abs(temp) % 10)

                # Iterates each digit across the row and then down the column and sets pixels_matrix
                for _ in range(8): # rows
                    for _ in range(4): # columns
                        # Update pixels_matrix image (pixels) from pixels model of each digit
                        self.pixels_matrix[index] = digits_representation[left_digit][digitIndex] # Left digit
                        self.pixels_matrix[index+4] = digits_representation[right_digit][digitIndex] # Right digit
                        index = index + 1 # Move to the next colum of the pixels_matrix
                        digitIndex = digitIndex + 1 # Move to the next pixel of the digit
                    index = index + 4 # Move to the next row of the pixels_matrix

                # If temperature < zero, add a minus before the digits
                if temp < 0:
                    self.pixels_matrix[24] = X
                else:
                    self.pixels_matrix[24] = O
            # Refresh the display
            if temp >= 30:
                if not self.blinking:
                    self.blinking = True
                    thread = Thread(target=self.blink)
                    thread.start()
            else:
                self.blinking = False
                self.sense.set_pixels(self.pixels_matrix)


    def blink(self):
        """Makes the display blink"""
        while self.temp >= 30:
            self.sense.set_pixels(self.pixels_matrix)
            sleep(0.4)
            self.sense.set_pixels(self.pixels_matrix_off)
            sleep(0.2)