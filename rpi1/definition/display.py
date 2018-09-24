from definition.colours import colour
from definition.digits import digits


def display(temp, ledDisplay):
    X, O = colour(temp)

    # Since the display is only big enough for two digits, an exception is made for anything over 99
    # If the temp does hit +/- 100 then it will blank out the display given that it still works
    if abs(temp) >= 100:
        for i in range(64):
            ledDisplay[i] = X

    else:	
        # Start building the display array
        # Iterates each digit across the row and then down the column
        index = 0
        digitIndex = 0
        digitBits = digits(X, O)

        for rowLoop in range(8):
                for columnLoop in range(4):
                        # Number 1 starts at position 32. Zero is 0 - 31, so multiply by 32
                        ledDisplay[index+4] = digitBits[int(abs(temp)%10)*32 + digitIndex]  # Second digit
                        ledDisplay[index] = digitBits[int(abs(temp)/10)*32 + digitIndex]  # First digit
                        index = index + 1
                        digitIndex = digitIndex + 1
                index = index + 4  # Move to the next row

        if temp < 0:
            ledDisplay[24] = X

