from models.digits import digits


def colour(temp):
    # Color code the display based on temperature range
    
    if(temp >= 40):
        X = [255, 0, 0]
    elif(temp >= 30):
        X = [255, 128, 0]
    elif(temp >= 20):
        X = [255, 255, 0]
    elif(temp >= 10):
        X = [0, 255, 0]
    elif(temp >= 0):
        X = [0, 255, 128]
    elif(temp >= -10):
        X = [0, 255, 255]
    elif(temp >= -20):
        X = [0, 191, 255]
    elif(temp >= -30):
        X = [0, 127, 255]
    elif(temp >= -40):
        X = [0, 64, 255]
    else:
        X = [0, 0, 255]

    O = [0, 0, 0]

    return X,O


def update_display_matrix(temp, ledDisplay):
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

