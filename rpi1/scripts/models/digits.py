def digits(X,O):
    """Returns the matrix representation of a digit, for the sense hat led panel"""
    # X are pixels (leds) that are ON
    # O are pixels (leds) that are OFF
    digits_representation = (
        # Zero
        O, X, X, X,
        O, X, O, X,
        O, X, O, X,
        O, X, O, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, O, O, O,

        # One
        O, O, X, O,
        O, X, X, O,
        O, O, X, O,
        O, O, X, O,
        O, O, X, O,
        O, O, X, O,
        O, X, X, X,
        O, O, O, O,

        # Two
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, X, X, X,
        O, X, O, O,
        O, X, O, O,
        O, X, X, X,
        O, O, O, O,

        # Three
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, X, X, X,
        O, O, O, O,

        # Four
        O, X, O, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, O,

        # Five
        O, X, X, X,
        O, X, O, O,
        O, X, O, O,
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, X, X, X,
        O, O, O, O,

        # Six
        O, X, X, X,
        O, X, O, O,
        O, X, O, O,
        O, X, X, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, O, O, O,

        # Seven
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, X,
        O, O, O, O,

        # Eight
        O, X, X, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, O, O, O,

        # Nine
        O, X, X, X,
        O, X, O, X,
        O, X, O, X,
        O, X, X, X,
        O, O, O, X,
        O, O, O, X,
        O, X, X, X,
        O, O, O, O
    )

    return digits_representation
