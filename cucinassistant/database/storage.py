from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user, check_number

from collections import namedtuple


Article = namedtuple('Article', ('aid', 'name', 'expiration', 'quantity'))


@use_user
def get_storage(cursor, uid):
    # Gets the contente
    cursor.execute('SELECT id, name, quantity, expiration FROM storage WHERE user=?;', [uid])
    return tuple(Article(*a) for a in cursor.fetchall())

@use_user
def get_storage_article(cursor, uid, aid):
    # Makes sure the id is valid
    if not (isinstance(aid, int) or isinstance(aid, str)) or (isinstance(aid, str) and not aid.isnumeric()):
        raise CAError('Elemento non valido')

    # Gets the contente
    cursor.execute('SELECT id, name, quantity, expiration FROM storage WHERE user=? AND id=?;', [uid, aid])
    if (data := cursor.fetchone()):
        return tuple(Article(*a) for a in cursor.fetchall())
    else:
        raise CAError('Articolo non trovato')

@use_user
def append_storage(cursor, uid, data_raw):
    data = []

    # Ensures the given data is valid
    if not data_raw: return
    for record in data_raw:
        if len(record) != 3 or not record[0]:
            raise CAError('Articolo non valido')
        elif record[2] and not check_number(record[2]):
            raise CAError('Quantit&agrave; non valida')
        elif record[1]:
            try:
                date.fromisoformat(record[1])
            except ValueError:
                raise CAError('Scadenza non valida')

        data.append((uid, record[0], record[1] or '2004-02-05', record[2] or '0'))

    # Adds the articles
    try:
        cursor.executemany('INSERT INTO storage (user, name, expiration, quantity) VALUES (?, ?, ?, ?);', data)
    except MDBError:
        raise CAError('Articolo gi&agrave; in lista')

@use_user
def remove_storage(cursor, uid, section, eids):
    # Cleans the list
    data = set()
    if not aids: return
    for aid in aids:
        if not aid:
            continue
        elif not check_number(aid):
            raise CAError('Articolo/i non valido/i')
        else:
            data.add((uid, aid))


    # Remove the selected items
    cursor.executemany('DELETE FROM storage WHERE user=? AND id=?;', list(data))
    if cursor.rowcount != len(data):
        raise CAError('Articolo/i non trovato/i')

# @use_user
# def edit_list(cursor, uid, section, eid, name):
#     if section not in ('shopping', 'ideas'): raise CAError('Lista inesistente')

#     # Ensures the items exists
#     if not (isinstance(eid, int) or isinstance(eid, str)) or (isinstance(eid, str) and not eid.isnumeric()):
#         raise CAError('Elemento non valido')
#     else:
#         cursor.execute(f'SELECT 1 FROM {section} WHERE id=? AND user=?;', [eid, uid])
#         if not cursor.fetchone():
#             raise CAError('Elemento non trovato')

#     # Continue only if the new name is new
#     cursor.execute(f'SELECT name FROM {section} WHERE id=?;', [eid])
#     if cursor.fetchone()[0] == name:
#         return

#     # Ensures the new name is valid
#     if not name:
#         raise CAError('Nuovo nome non valido')

#     # Makes sure tha name is unique
#     cursor.execute(f'SELECT 1 FROM {section} WHERE name=? AND user=?;', [name, uid])
#     if cursor.fetchone():
#         raise CAError('Elemento gi&agrave; in lista')

#     # Saves the change
#     cursor.execute(f'UPDATE {section} SET name=? WHERE id=?;', [name, eid])
