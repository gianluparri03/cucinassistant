from . import CAError

from hashlib import sha256
from functools import wraps
from secrets import token_hex
from sqlite3 import connect, IntegrityError


def hash_password(password):
    return sha256(password.encode('utf-8')).hexdigest()

def init_db():
    global db
    db = connect('cucinassistant.db', check_same_thread=False, isolation_level=None)

    db.execute('''CREATE TABLE IF NOT EXISTS users (
                  uid INTEGER PRIMARY KEY AUTOINCREMENT,
                  username TEXT NOT NULL UNIQUE,
                  password TEXT NOT NULL,
                  email TEXT NOT NULL UNIQUE,
                  token TEXT,
                  newsletter BOOLEAN DEFAULT TRUE);''')

    db.execute('''CREATE TABLE IF NOT EXISTS menus (
                  user INTEGER PRIMARY KEY REFERENCES users (uid),
                  menu TEXT);''')

    db.execute('''CREATE TABLE IF NOT EXISTS storage (
                  user INTEGER REFERENCES users (uid),
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  name TEXT NOT NULL,
                  quantity INT,
                  expiration DATE,
                  UNIQUE (user, name, expiration));''')

    db.execute('''CREATE TABLE IF NOT EXISTS shopping (
                  user INTEGER REFERENCES users (uid),
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  name TEXT NOT NULL,
                  UNIQUE (user, name));''')

    db.execute('''CREATE TABLE IF NOT EXISTS ideas (
                  user INTEGER REFERENCES users (uid),
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  name TEXT NOT NULL,
                  UNIQUE (user, name));''')

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
        cursor.execute('INSERT INTO users (username, email, password) VALUES (?, ?, ?);', [username, email, password])
        return cursor.lastrowid

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
    cursor.execute('SELECT uid FROM users WHERE username=? AND password=?;', [username, password])
    if not (uid := cursor.fetchone()):
        raise CAError('Credenziali non valide')
    else:
        return uid[0]


@use_db
def get_user_username(cursor, uid):
    # Returns the user's username
    cursor.execute('SELECT username FROM users WHERE uid=?;', [uid])
    if (username := cursor.fetchone()):
        return username[0]

@use_db
def get_user_email(cursor, uid):
    # Returns the user's email
    cursor.execute('SELECT email FROM users WHERE uid=?;', [uid])
    if (email := cursor.fetchone()):
        return email[0]


@use_db
def generate_user_token(cursor, uid):
    # Generates a new deletion token for the user
    token = token_hex(18)
    cursor.execute('UPDATE users SET token=? WHERE uid=?;', [token, uid])
    return token

@use_db
def delete_user(cursor, uid, token):
    # Checks if the token is valid
    cursor.execute('SELECT 1 FROM users WHERE uid=? AND token=?;', [uid, token])
    if not cursor.fetchone():
        raise CAError('Errore durante la cancellazione, riprova.')

    cursor.execute('DELETE FROM users WHERE uid=?;', [uid])


@use_db
def change_user_email(cursor, uid, new):
    # Saves the new one
    cursor.execute('UPDATE users SET email=? WHERE uid=?;', [new, uid])

@use_db
def change_user_password(cursor, uid, old, new):
    old = hash_password(old)
    new = hash_password(new)

    # Checks if the old are valid
    cursor.execute('SELECT 1 FROM users WHERE uid=? AND password=?;', [uid, old])
    if not cursor.fetchone():
        raise CAError('Password attuale non valida')

    # Saves the new ones
    cursor.execute('UPDATE users SET password=? WHERE uid=?;', [new, uid])

@use_db
def reset_user_password(cursor, email):
    # Selects the uid
    cursor.execute('SELECT username FROM users WHERE email=?;', [email])
    username = cursor.fetchone()
    if not username:
        return

    # Generates a new password and saves it 
    unhashed = token_hex(4)
    hashed = hash_password(unhashed)
    cursor.execute('UPDATE users SET password=? WHERE username=?;', [hashed, username[0]])

    return username[0], unhashed


@use_db
def get_user_menu(cursor, uid):
    # Returns the menu
    cursor.execute('SELECT menu FROM menus WHERE user=?;', [uid])
    if (menu := cursor.fetchone()):
        return menu[0].split(';')
    else:
        return [] * 14

@use_db
def update_user_menu(cursor, uid, items):
    # Saves the new menu
    if len(items) != 14:
        raise CAError('Menu non valido')
    cursor.execute('REPLACE INTO menus (user, menu) VALUES (?, ?);', [uid, ';'.join(items)])


@use_db
def get_user_storage(cursor, uid):
    # Returns the storage
    cursor.execute('SELECT id, name, quantity, expiration FROM storage WHERE user=? ORDER BY expiration;', [uid])
    return cursor.fetchall()

@use_db
def add_user_storage(cursor, uid, items):
    # Adds the items, or updates them
    for item in items:
        try:
            cursor.execute('INSERT INTO storage (user, name, quantity, expiration) VALUES (?, ?, ?, ?);', [uid, item[0], item[1], item[2]])
        except IntegrityError:
            cursor.execute('UPDATE storage SET quantity=quantity+? WHERE user=? AND name=? AND expiration=?;', [item[1], uid, item[0], item[2]])

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
    new = max(eval(str(data[1] or '0') + delta), 0)
    if new > 0:
        cursor.execute('UPDATE storage SET quantity=? WHERE id=?;', [new, item])
    else:
        cursor.execute('DELETE FROM storage WHERE id=?;', [item])

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
    return {l[0]: l[1] for l in cursor.fetchall()}

@use_db
def add_user_lists(cursor, uid, section, items):
    if section not in ('shopping', 'ideas'): raise CAError('Sezione sconosciuta')

    # Adds some items to the list
    data = [[uid, item] for item in items if item]
    cursor.executemany(f'INSERT OR IGNORE INTO {section} (user, name) VALUES (?, ?);', data)

@use_db
def remove_user_lists(cursor, uid, section, items):
    if section not in ('shopping', 'ideas'): raise CAError('Sezione sconosciuta')

    # Remove some items from the list
    data = [[uid, item] for item in items if item]
    cursor.executemany(f'DELETE FROM {section} WHERE user=? AND id=?;', data)


@use_db
def get_users_number(cursor):
    # Counts the users
    cursor.execute('SELECT COUNT(*) FROM users;')
    return cursor.fetchone()[0]

@use_db
def get_newsletter_emails(cursor):
    # Returns the emails of the users that receives the newsletter
    cursor.execute('SELECT email FROM users WHERE newsletter;')
    return [c[0] for c in cursor.fetchall() if c]

@use_db
def disable_newsletter(cursor, uid):
    # Disables the user's newsletter
    cursor.execute('UPDATE users SET newsletter=0 WHERE uid=?;', [uid])

@use_db
def enable_newsletter(cursor, uid):
    # Enables the user's newsletter
    cursor.execute('UPDATE users SET newsletter=1 WHERE uid=?;', [uid])
