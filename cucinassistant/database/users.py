from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_db, ph

from secrets import token_hex
from mariadb import Error as MDBError
from string import ascii_letters, digits
from argon2.exceptions import VerificationError


@use_db
def get_users_number(cursor):
    # Counts the users
    cursor.execute('SELECT COUNT(*) FROM users;')
    return cursor.fetchone()[0]

@use_db
def get_users_emails(cursor):
    # Counts the users
    cursor.execute('SELECT email FROM users;')
    return list(map(lambda r: r[0], cursor.fetchall()))

@use_db
def create_user(cursor, username, email, password):
    # Makes some checks
    if len(username) < 3:
        raise CAError('Nome utente non valido (lunghezza minima 3 caratteri)')
    elif set(username) - set(ascii_letters + digits + '_'):
        raise CAError('Nome utente non valido (solo lettere, numeri e "_" consentiti)')
    elif len(password) < 5:
        raise CAError('Password non valida (lunghezza minima 5 caratteri)')

    try:
        # Tries to create a new user
        password = ph.hash(password)
        cursor.execute('INSERT INTO users (username, email, password) VALUES (?, ?, ?);', [username, email, password])
        return cursor.lastrowid

    # Rewrites the error
    except MDBError as e:
        if str(e).endswith("for key 'email'"):
            raise CAError("Email non disponibile")
        elif str(e).endswith("for key 'username'"):
            raise CAError("Nome utente non disponibile")
        else:
            raise CAError("Errore sconosciuto")

@use_db
def login_user(cursor, username, password):
    # Checks if the credentials are valid
    cursor.execute('SELECT uid, password FROM users WHERE username=?;', [username])
    try:
        if (data := cursor.fetchone()):
            ph.verify(data[1], password)
            return data[0]
        else:
            raise VerificationError()
    except VerificationError:
        raise CAError('Credenziali non valide')

@use_db
def get_user_data(cursor, uid, email=''):
    # Returns the user's data
    if email:
        cursor.execute('SELECT uid, username, email, password, token FROM users WHERE email=?;', [email])
    else:
        cursor.execute('SELECT uid, username, email, password, token FROM users WHERE uid=?;', [uid])

    if (data := cursor.fetchone()):
        return dict(zip(('uid', 'username', 'email', 'password', 'token'), data))
    else:
        raise CAError('Utente sconosciuto')

@use_db
def generate_user_token(cursor, uid):
    # Generates a new deletion token for the user
    token = token_hex(18)
    cursor.execute('UPDATE users SET token=? WHERE uid=?;', [ph.hash(token), uid])
    if cursor.rowcount == 1:
        return token
    else:
        raise CAError('Utente sconosciuto')

@use_db
def delete_user(cursor, uid, token):
    # Checks if the token is valid
    cursor.execute('SELECT token FROM users WHERE uid=?;', [uid])
    try:
        if (data := cursor.fetchone()):
            ph.verify(data[0], token)
            cursor.execute('DELETE FROM users WHERE uid=?;', [uid])
        else:
            raise VerificationError()
    except VerificationError:
        raise CAError('Errore durante la cancellazione, riprova')

@use_db
def change_user_username(cursor, uid, new):
    # Saves the new one
    try:
        if not (data := get_user_data(uid)):
            raise CAError('Utente sconosciuto')
        elif data.get('username') == new:
            return

        cursor.execute('UPDATE users SET username=? WHERE uid=?;', [new, uid])
    except MDBError:
        raise CAError('Nome utente non disponibile')

@use_db
def change_user_email(cursor, uid, new):
    # Saves the new one
    try:
        if not (data := get_user_data(uid)):
            raise CAError('Utente sconosciuto')
        elif data.get('email') == new:
            return

        cursor.execute('UPDATE users SET email=? WHERE uid=?;', [new, uid])
    except MDBError:
        raise CAError('Email non disponibile')

@use_db
def change_user_password(cursor, uid, old, new):
    # Ensures that the user exists
    if not get_user_data(uid):
        raise CAError('Utente sconosciuto')

    cursor.execute('SELECT password FROM users WHERE uid=?;', [uid])
    try:
        # Check if the user is athorized, then updates it
        ph.verify(cursor.fetchone()[0], old)
        cursor.execute('UPDATE users SET password=? WHERE uid=?;', [ph.hash(new), uid])
    except VerificationError:
        raise CAError('Credenziali non valide')

@use_db
def reset_user_password(cursor, email, token, new):
    # Ensures that the user exists
    data = get_user_data('', email=email)
    if not data:
        raise CAError('Utente sconosciuto')

    try:
        # Check if the user is athorized, then updates it
        ph.verify(data['token'], token)
        cursor.execute('UPDATE users SET password=? WHERE uid=?;', [ph.hash(new), data['uid']])
    except VerificiationError:
        raise CAError('Errore durante la reimpostazione della password')
