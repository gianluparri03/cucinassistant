from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user, check_number

from collections import namedtuple
from functools import wraps


Entry = namedtuple('Entry', ('eid', 'name'))


def check_section(func):
    @wraps(func)
    @use_user
    def inner(cursor, uid, section, *args, **kwargs):
        # Ensures the section exists
        if section not in ('shopping', 'ideas'):
            raise CAError('Lista inesistente')

        return func(cursor, uid, section, *args, **kwargs)

    return inner


@check_section
def get_list(cursor, uid, section):
    # Gets the list
    cursor.execute(f'SELECT id, name FROM {section} WHERE user=?;', [uid])
    return tuple(Entry(l[0], l[1]) for l in cursor.fetchall())

@check_section
def get_list_entry(cursor, uid, section, eid):
    # Makes sure the id is valid
    if not check_number(eid):
        raise CAError('Elemento non valido')

    # Gets the entry
    cursor.execute(f'SELECT id, name FROM {section} WHERE user=? AND id=?;', [uid, int(eid)])
    if (n := cursor.fetchone()):
        return Entry(*n)
    else:
        raise CAError('Elemento non in lista')

@check_section
def append_list(cursor, uid, section, names):
    # Appends the items to the list
    data = [[uid, name] for name in names if name]
    if not data: return
    cursor.executemany(f'INSERT IGNORE INTO {section} (user, name) VALUES (?, ?);', data)

@check_section
def remove_list(cursor, uid, section, eids):
    # Cleans the list
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

    # Remove some items from the list
    cursor.executemany(f'DELETE FROM {section} WHERE user=? AND id=?;', list(data))
    if cursor.rowcount != len(data):
        raise CAError('Elemento non trovato')

@check_section
def edit_list(cursor, uid, section, eid, name):
    # Ensures the items exists
    if not check_number(eid):
        raise CAError('Elemento non valido')
    else:
        cursor.execute(f'SELECT name FROM {section} WHERE id=? AND user=?;', [eid, uid])
        if not (data := cursor.fetchone()):
            raise CAError('Elemento non trovato')

        # Continue only if the new name is new
        elif data[0] == name:
            return

    # Ensures the new name is valid
    if not name:
        raise CAError('Nuovo nome non valido')

    # Makes sure tha name is unique
    cursor.execute(f'SELECT 1 FROM {section} WHERE name=? AND user=?;', [name, uid])
    if cursor.fetchone():
        raise CAError('Elemento gi&agrave; in lista')

    # Saves the change
    cursor.execute(f'UPDATE {section} SET name=? WHERE id=?;', [name, eid])
