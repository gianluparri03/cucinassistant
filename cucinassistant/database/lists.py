from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user

from collections import namedtuple


Entry = namedtuple('Entry', ('eid', 'name'))


@use_user
def get_list(cursor, uid, section):
    if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

    # Gets the list
    cursor.execute(f'SELECT id, name FROM {section} WHERE user=?;', [uid])
    return tuple(Entry(l[0], l[1]) for l in cursor.fetchall())

@use_user
def get_list_entry(cursor, uid, section, eid):
    if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

    # Makes sure the id is valid
    if not (isinstance(eid, int) or isinstance(eid, str)) or (isinstance(eid, str) and not eid.isnumeric()):
        raise CAError('Elemento non valido')

    # Gets the entry
    cursor.execute(f'SELECT name FROM {section} WHERE user=? AND id=?;', [uid, eid])
    if (n := cursor.fetchone()):
        return n[0]
    else:
        raise CAError('Articolo non in lista')

@use_user
def append_list(cursor, uid, section, names):
    if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

    # Appends the items to the list
    data = [[uid, name] for name in names if name]
    if not data: return
    cursor.executemany(f'INSERT IGNORE INTO {section} (user, name) VALUES (?, ?);', data)
    return cursor.lastrowid

@use_user
def remove_list(cursor, uid, section, eids):
    if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

    # Cleans the list
    data = set()
    for eid in eids:
        if not eid:
            continue
        elif not (isinstance(eid, int) or isinstance(eid, str)) or (isinstance(eid, str) and not eid.isnumeric()):
            raise CAError('Elemento/i non valido/i')
        else:
            data.add((uid, eid))

    if not data:
        return

    # Remove some items from the list
    cursor.executemany(f'DELETE FROM {section} WHERE user=? AND id=?;', list(data))
    if cursor.rowcount != len(data):
        raise CAError('Elemento/i non trovato/i')

@use_user
def edit_list(cursor, uid, section, eid, name):
    if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

    # Ensures the items exists
    if not (isinstance(eid, int) or isinstance(eid, str)) or (isinstance(eid, str) and not eid.isnumeric()):
        raise CAError('Elemento non valido')
    else:
        cursor.execute(f'SELECT 1 FROM {section} WHERE id=? AND user=?;', [eid, uid])
        if not cursor.fetchone():
            raise CAError('Elemento non trovato')

    # Ensures the new name doesn't exist
    if not name:
        raise CAError('Nuovo nome non valido')
    else:
        cursor.execute(f'SELECT 1 FROM {section} WHERE name=? AND user=?;', [name, uid])
        if cursor.fetchone():
            raise CAError('Elemento gi&agrave; in lista')

    # Saves the change
    cursor.execute(f'UPDATE {section} SET name=? WHERE id=?;', [name, eid])
