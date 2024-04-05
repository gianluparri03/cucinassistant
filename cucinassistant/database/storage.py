from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user, check_number

from mariadb import Error as MDBError
from collections import namedtuple
from datetime import date


Article = namedtuple('Article', ('aid', 'name', 'expiration', 'quantity'))

NULL_EXPIRATION = '2004-02-05'
NULL_QUANTITY = '0'

def new_article(a, n, e, q):
    return Article(a, n, e.isoformat() if e.isoformat() != NULL_EXPIRATION else None, q if str(q) != NULL_QUANTITY else None)

def verify_article(data):
    if len(data) != 3:
        raise CAError('Formato non valido')
    elif not data[0]:
        raise CAError('Nome non valido')
    elif data[2] and not check_number(data[2]):
        raise CAError('Quantit&agrave; non valida')
    elif data[1]:
        try:
            date.fromisoformat(data[1])
        except ValueError:
            raise CAError('Scadenza non valida')

    return data[0], data[1] or NULL_EXPIRATION, data[2] or NULL_QUANTITY


@use_user
def get_storage(cursor, uid, name=''):
    # Gets the contente
    cursor.execute(f'SELECT id, name, expiration, quantity FROM storage WHERE user=? AND name LIKE ? ORDER BY expiration;', [uid, '%' + (name or '') + '%'])
    return tuple(new_article(*a) for a in cursor.fetchall())

@use_user
def get_storage_article(cursor, uid, aid):
    # Makes sure the id is valid
    if not (isinstance(aid, int) or isinstance(aid, str)) or (isinstance(aid, str) and not aid.isnumeric()):
        raise CAError('Articolo non valido')

    # Gets the contente
    cursor.execute('SELECT id, name, expiration, quantity FROM storage WHERE user=? AND id=?;', [uid, aid])
    if (data := cursor.fetchone()):
        return new_article(*data)
    else:
        raise CAError('Articolo non trovato')

@use_user
def append_storage(cursor, uid, data_raw):
    data = []

    # Ensures the given data is valid
    if not data_raw: return
    for record in data_raw:
        data.append((uid, *verify_article(record)))

    # Adds the articles
    try:
        cursor.executemany('INSERT INTO storage (user, name, expiration, quantity) VALUES (?, ?, ?, ?);', data)
    except MDBError:
        raise CAError('Articolo gi&agrave; presente')

@use_user
def remove_storage(cursor, uid, aids):
    # Cleans the list
    data = set()
    for aid in aids:
        if not aid:
            continue
        elif not check_number(aid):
            raise CAError('Articolo non valido')
        else:
            data.add((uid, aid))
    if not data: return

    # Remove the selected items
    cursor.executemany('DELETE FROM storage WHERE user=? AND id=?;', list(data))
    if cursor.rowcount != len(data):
        raise CAError('Articolo non trovato')

@use_user
def edit_storage(cursor, uid, aid, data):
    # Ensures the items exists
    if not check_number(aid):
        raise CAError('Articolo non valido')
    else:
        cursor.execute('SELECT name, expiration, quantity FROM storage WHERE id=? AND user=?;', [aid, uid])
        record = cursor.fetchone()
        if not data:
            raise CAError('Articolo non trovato')

        # Verifies the article
        name, exp, qty = verify_article(data)

        # Continue only if something is new
        if name == record[0] and exp == record[1].isoformat() and qty == str(record[2]):
            return

    # Makes sure tha name is unique
    cursor.execute('SELECT id FROM storage WHERE name=? AND expiration=? AND id != ? AND user=?;', [name, exp, aid, uid])
    if cursor.fetchone():
        raise CAError('Articolo gi&agrave; presente')

    # Saves the change
    cursor.execute('UPDATE storage SET name=?, expiration=?, quantity=? WHERE id=?;', [name, exp, qty, aid])
