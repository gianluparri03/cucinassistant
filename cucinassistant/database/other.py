from cucinassistant.exceptions import CAError
from cucinassistant.database import db, ph, use_db

from datetime import date
from functools import wraps
from secrets import token_hex
from argon2 import PasswordHasher
from string import ascii_letters, digits
from mariadb import connect, Error as MDBError
from argon2.exceptions import VerificationError

# TODO split into multiple files



@use_db
def get_user_storage(cursor, uid):
    # Returns the storage
    cursor.execute("SELECT id, name, quantity, DATE_FORMAT(expiration, '%d\\/%m\\/%Y') FROM storage WHERE user=? ORDER BY expiration;", [uid])
    if (data := cursor.fetchall()):
        return [[i[0], i[1], i[2] if i[2] != 0 else '', i[3] if not i[3] == '05/02/2004' else ''] for i in data]
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
    except MDBError:
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
