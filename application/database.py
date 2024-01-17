from . import db

from hashlib import sha256
from functools import wraps
from secrets import token_hex
from sqlite3 import IntegrityError


class CAError(Exception): pass

def hash_password(password):
    return sha256(password.encode('utf-8')).hexdigest()

def use_db(func):
    @wraps(func)
    def inner(*args, **kwargs):
        cur = db.cursor()
        ris = func(cur, *args, **kwargs)
        cur.close()
        return ris

    return inner


@use_db
def create_user(cursor, username, email, password):
    try:
        # Tries to create a new user
        password = hash_password(password)
        cursor.execute('INSERT INTO users (username, email, password) VALUES (?, ?, ?);', (username, email, password))

    # Rewrites the error
    except IntegrityError as e:
        match str(e):
            case 'UNIQUE constraint failed: users.email':
                raise CAError("Email non disponibile")
            case 'UNIQUE constraint failed: users.username':
                raise CAError("Username non disponibile")
            case _:
                raise CAError("Errore sconosciuto")

@use_db
def login_user(cursor, username, password):
    # Checks if the credentials are valid
    password = hash_password(password)
    cursor.execute('SELECT 1 FROM users WHERE username=? AND password=?;', (username, password))
    if not cursor.fetchone():
        raise CAError('Credenziali non valide')


@use_db
def get_user_email(cursor, username):
    # Returns the user's email
    cursor.execute('SELECT email FROM users WHERE username=?;', (username, ))
    if (res := cursor.fetchone()):
        return res[0]


@use_db
def generate_user_token(cursor, username):
    # Generates a new deletion token for the user
    token = token_hex(18)
    cursor.execute('UPDATE users SET token=? WHERE username=?;', (token, username))
    return token

@use_db
def delete_user(cursor, username, token):
    # Checks if the token is valid
    cursor.execute('SELECT 1 FROM users WHERE username=? AND token=?;', (username, token))
    if not cursor.fetchone():
        raise CAError('Errore durante la cancellazione, riprova.')

    cursor.execute('DELETE FROM users WHERE username=?;', (username, ))


@use_db
def change_user_password(cursor, username, old, new):
    old = hash_password(old)
    new = hash_password(new)

    # Checks if the old are valid
    cursor.execute('SELECT 1 FROM users WHERE username=? AND password=?;', (username, old))
    if not cursor.fetchone():
        raise CAError('Password attuale non valida')

    # Saves the new ones
    cursor.execute('UPDATE users SET password=? WHERE username=?;', (new, username))

@use_db
def reset_user_password(cursor, email):
    # Selects the username
    cursor.execute('SELECT username FROM users WHERE email=?;', (email, ))
    data = cursor.fetchone()
    if not data:
        return

    # Generates a new password and saves it 
    unhashed = token_hex(4)
    hashed = hash_password(unhashed)
    cursor.execute('UPDATE users SET password=? WHERE username=?;', (hashed, data[0]))

    return data[0], unhashed
