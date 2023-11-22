from json import load, dump
from hashlib import sha256
from os import remove as os_remove


BASE_FOLDER = 'data/'
USERS_FILE  = '_users'

def read_file(filename):
    # Reads the given file
    try:
        with open(BASE_FOLDER + filename + '.json') as f:
            return load(f)
    except FileNotFoundError:
        return {}

def write_file(filename, content):
    # Writes the given content in the given file
    with open(BASE_FOLDER + filename + '.json', 'w') as f:
        dump(content, f)


class User:
    def __init__(self, username):
        self.username = username

    @staticmethod
    def hash_password(password):
        h = sha256(password.encode('utf-8'))
        return h.hexdigest()

    @staticmethod # Returns an error
    def create(username, email, password):
        users = read_file(USERS_FILE)

        # Returns an error if not valid
        if username in users:
            return 'Nome utente non disponibile'
        elif not username.isalnum() or not password.isalnum():
            return 'Nome utente e password devono essere alfanumerici'

        # Creates the user
        users[username] = [email, User.hash_password(password)]

        # Creates the files
        write_file(USERS_FILE, users)
        write_file(username, {"menu": [""] * 14, "spesa": [], "idee": []})

    @staticmethod # Returns an error
    def check(username, password):
        # Checks if the credentials are valid
        data = read_file(USERS_FILE).get(username)
        if not data or User.hash_password(password) != data[1]:
            return 'Credenziali non valide'

    def delete(self):
        # Deletes the user
        write_file(USERS_FILE, {k: v for k, v in read_file(USERS_FILE).items() if k != self.username})
        try:
            os_remove(BASE_FOLDER + self.username + '.json')
        except OSError:
            pass

    def read_data(self, part=""):
        # Reads the user's data
        data = read_file(self.username)
        return data[part] if part else data

    def update_data(self, part_name="", part_data=None):
        # Updates the user's data
        data = self.read_data() | {part_name: part_data}
        write_file(self.username, data)
