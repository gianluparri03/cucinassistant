from uuid import uuid4
from hashlib import sha256
from json import load, dump
from os import remove as os_remove


BASE_FOLDER = 'data/'
USERS_FILE  = '_users'
get_filename = lambda f: BASE_FOLDER + f + '.json'

def read_file(filename):
    # Reads the given file
    try:
        with open(get_filename(filename)) as f:
            return load(f)
    except FileNotFoundError:
        return {}

def write_file(filename, content):
    # Writes the given content in the given file
    with open(get_filename(filename), 'w') as f:
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
        elif email in [u[0] for u in users.values()]:
            return 'Email non disponibile'
        elif not username.isalnum() or not password.isalnum():
            return 'Nome utente e password devono essere alfanumerici'

        # Creates the user
        users[username] = [email, User.hash_password(password)]

        # Creates the files
        write_file(USERS_FILE, users)
        write_file(username, {"menu": [""] * 14, "spesa": [], "idee": [], "scadenze": [], "quantita": []})

    @staticmethod # Returns an error
    def check(username, password):
        # Checks if the credentials are valid
        data = read_file(USERS_FILE).get(username)
        if not data or User.hash_password(password) != data[1]:
            return 'Credenziali non valide'

    @staticmethod # Returns the username and the new password, or None
    def reset_password(email):
        users = read_file(USERS_FILE)
        username = ''

        # Looks for the username
        for u in users:
            if users[u][0] == email:
                username = u
        if not username: return

        # Saves the new password
        password = str(uuid4()).split('-')[0]
        users[username][1] = User.hash_password(password)
        write_file(USERS_FILE, users)

        return username, password

    @staticmethod
    def get_number():
        return len(read_file(USERS_FILE))

    def delete(self):
        # Deletes the user
        write_file(USERS_FILE, {k: v for k, v in read_file(USERS_FILE).items() if k != self.username})
        try:
            os_remove(get_filename(self.username))
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
