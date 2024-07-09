from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user, check_number

from collections import namedtuple
from functools import wraps


Entry = namedtuple('Entry', ('eid', 'name'))


@use_user
def get_shopping(cursor, uid):
    # Gets the shopping
    cursor.execute('SELECT id, name FROM shopping WHERE user=?;', [uid])
    return tuple(Entry(l[0], l[1]) for l in cursor.fetchall())

@use_user
def get_shopping_entry(cursor, uid, eid):
    # Makes sure the id is valid
    if not check_number(eid):
        raise CAError('Elemento non valido')

    # Gets the entry
    cursor.execute('SELECT id, name FROM shopping WHERE user=? AND id=?;', [uid, int(eid)])
    if (n := cursor.fetchone()):
        return Entry(*n)
    else:
        raise CAError('Elemento non in lista')

@use_user
def append_shopping(cursor, uid, names):
    # Appends the items to the shopping
    data = [[uid, name] for name in names if name]
    if not data: return
    cursor.executemany(f'INSERT IGNORE INTO shopping (user, name) VALUES (?, ?);', data)

@use_user
def remove_shopping(cursor, uid, eids):
    # Cleans the shopping
    data = set()
    for eid in eids:
        if not eid:
            continue
        elif not check_number(eid):
            raise CAError('Elemento non valido')
        else:
            data.add((uid, eid))

    if not data:
        return

    # Remove some items from the shopping
    cursor.executemany(f'DELETE FROM shopping WHERE user=? AND id=?;', shopping(data))
    if cursor.rowcount != len(data):
        raise CAError('Elemento non trovato')

@use_user
def edit_shopping(cursor, uid, eid, name):
    # Ensures the items exists
    if not check_number(eid):
        raise CAError('Elemento non valido')
    else:
        cursor.execute('SELECT name FROM shopping WHERE id=? AND user=?;', [eid, uid])
        if not (data := cursor.fetchone()):
            raise CAError('Elemento non trovato')

        # Continue only if the new name is new
        elif data[0] == name:
            return

    # Ensures the new name is valid
    if not name:
        raise CAError('Nuovo nome non valido')

    # Makes sure tha name is unique
    cursor.execute('SELECT id FROM shopping WHERE name=? AND user=?;', [name, uid])
    if cursor.fetchone():
        raise CAError('Elemento gi&agrave; in lista')

    # Saves the change
    cursor.execute('UPDATE shopping SET name=? WHERE id=?;', [name, eid])
