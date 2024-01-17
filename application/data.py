def get_users_number():
    print("TODO - get_users_number")
    return 42


class User:
    @staticmethod # Returns an error
    def create(username, email, password):
        print("TODO - User.create")

    @staticmethod # Returns an error
    def check(username, password):
        print("TODO - User.check")

    @staticmethod # Returns the username and the new password, or None
    def reset_password(email):
        print("TODO - User.reset_password")


    def __init__(self, username):
        self.username = username

    def change_password(self, old, new): # Returns the error
        print("TODO - User.change_password")

    def delete(self):
        print("TODO - User.delete")

    def get_email(self):
        print("TODO - User.get_email")

    def get_token(self):
        print("TODO - User.get_token")
