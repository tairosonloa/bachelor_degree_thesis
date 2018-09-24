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
