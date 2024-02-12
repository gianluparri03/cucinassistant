from . import config, CAError

import string
from datetime import date
from functools import wraps
from secrets import token_hex
from argon2 import PasswordHasher
from mariadb import connect, Error as DBError
from argon2.exceptions import VerificationError


def init_db():
    global db, ph
    db = connect(host=config['Database']['Hostname'], database=config['Database']['Database'], \
                 user=config['Database']['Username'], password=config['Database']['Password'])
    db.autocommit = True

    with db.cursor() as cursor:
        with open('application/schema.sql') as f:
            for command in f.read().split(';')[:-1]:
                cursor.execute(command)

    ph = PasswordHasher()

def use_db(func):
    @wraps(func)
    def inner(*args, **kwargs):
        db.reconnect()
        with db.cursor() as cur:
            return func(cur, *args, **kwargs)

    return inner


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
    elif set(username) - set(string.ascii_letters + string.digits + '_'):
        raise CAError('Nome utente non valido (solo lettere, numeri e "_" consentiti)')
    elif len(password) < 5:
        raise CAError('Password non valida (lunghezza minima 5 caratteri)')

    try:
        # Tries to create a new user
        password = ph.hash(password)
        cursor.execute('INSERT INTO users (username, email, password) VALUES (?, ?, ?);', [username, email, password])
        return cursor.lastrowid

    # Rewrites the error
    except DBError as e:
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
    except DBError:
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
    except DBError:
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


@use_db
def get_user_menu(cursor, uid):
    # Returns the menu
    cursor.execute('SELECT menu FROM menus WHERE user=?;', [uid])
    if (menu := cursor.fetchone()):
        return menu[0].split(';')
    else:
        # Ensures that the user exists
        cursor.execute('SELECT 1 FROM users WHERE uid=?;', [uid])
        if (data := cursor.fetchone()):
            return [] * 14
        else:
            raise CAError('Utente sconosciuto')

@use_db
def update_user_menu(cursor, uid, items):
    # Checks the menu syntax
    if len(items) != 14:
        raise CAError('Menu non valido')

    # Saves the new menu
    try:
        cursor.execute('REPLACE INTO menus (user, menu) VALUES (?, ?);', [uid, ';'.join(items)])
    except DBError:
        raise CAError('Utente sconosciuto')


@use_db
def get_user_storage(cursor, uid):
    # Returns the storage
    cursor.execute('SELECT id, name, quantity, expiration FROM storage WHERE user=? ORDER BY expiration;', [uid])
    if (data := cursor.fetchall()):
        return [[i[0], i[1], i[2] if i[2] != 0 else '', i[3] if not i[3].isoformat().startswith('2004-02-05') else ''] for i in data]
    else:
        # Ensures that the user exists
        cursor.execute('SELECT 1 FROM users WHERE uid=?;', [uid])
        if (data := cursor.fetchone()):
            return []
        else:
            raise CAError('Utente sconosciuto')


@use_db
def add_user_storage(cursor, uid, items):
    # Makes sure the items syntax is correct
    for item in items:
        if len(item) != 3:
            raise CAError('Formato non valido')
        if not item[0]:
            raise CAError('Nome non valido')
        if not item[1]:
            item[1] = 0
        elif not item[1].isnumeric():
            raise CAError('Quantità non valida')
        if not item[2]:
            item[2] = '2004-02-05'
        else:
            try:
                date.fromisoformat(item[2])
            except ValueError:
                raise CAError('Data di scandenza non valida')

    # Adds the items, or updates them
    try:
        cursor.executemany('INSERT INTO storage (user, name, quantity, expiration) VALUES (?, ?, ?, ?)' + \
            'ON DUPLICATE KEY UPDATE quantity=quantity+?;', [[uid, i[0], i[1], i[2], i[1]] for i in items])
    except DBError:
        raise CAError('Utente sconosciuto')

@use_db
def edit_user_storage(cursor, uid, item, delta):
    if not delta[1:].isnumeric() or (delta[0] != '-' and delta[0] != '+'):
        raise CAError('Valore non valido')

    # Fetches the old quantity
    cursor.execute('SELECT 1, quantity FROM storage WHERE user=? AND id=?;', [uid, item])
    data = cursor.fetchone()
    if not data:
        raise CAError('Articolo non in lista')

    # Save the changes
    new = max(eval(str(data[1]) + delta), 0)
    cursor.execute('UPDATE storage SET quantity=? WHERE id=?;', [new, item])

@use_db
def remove_user_storage(cursor, uid, items):
    # Remove some items from the storage
    data = [[uid, item] for item in items if item]
    cursor.executemany('DELETE FROM storage WHERE user=? AND id=?;', data)


@use_db
def get_user_lists(cursor, uid, section):
    if section not in ('shopping', 'ideas'): raise CAError('Sezione sconosciuta')

    # Returns the list
    cursor.execute(f'SELECT id, name FROM {section} WHERE user=?;', [uid])
    if (data := cursor.fetchall()):
        return {i[0]: i[1] for i in data}
    else:
        # Ensures that the user exists
        cursor.execute('SELECT 1 FROM users WHERE uid=?;', [uid])
        if (data := cursor.fetchone()):
            return {}
        else:
            raise CAError('Utente sconosciuto')

@use_db
def add_user_lists(cursor, uid, section, items):
    if section not in ('shopping', 'ideas'): raise CAError('Sezione sconosciuta')

    # Adds some items to the list
    try:
        data = [[uid, item] for item in items if item]
        cursor.executemany(f'INSERT IGNORE INTO {section} (user, name) VALUES (?, ?);', data)
    except DBError:
        raise CAError('Utente sconosciuto')

@use_db
def remove_user_lists(cursor, uid, section, items):
    if section not in ('shopping', 'ideas'): raise CAError('Sezione sconosciuta')

    # Remove some items from the list
    data = [[uid, item] for item in items if item]
    cursor.executemany(f'DELETE FROM {section} WHERE user=? AND id=?;', data)
