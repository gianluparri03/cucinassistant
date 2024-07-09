class CAError(Exception):
    # The general one
    pass

class CACritical(Exception):
    # Only for when the user has to be logged out
    # from the website
    pass
