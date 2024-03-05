from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_db

from collections import namedtuple
from mariadb import connect, Error as MDBError


Menu = namedtuple('Menu', ['menu', 'prev', 'next'])

@use_db
def get_user_menu(cursor, uid, mid=0):
    # Ensures that the user exists
    cursor.execute('SELECT 1 FROM users WHERE uid=?;', [uid])
    if not cursor.fetchone():
        raise CAError('Utente sconosciuto')

    menu = []

    if not mid:
        # Returns the first menu
        cursor.execute('SELECT menu FROM menus WHERE user=? ORDER BY id LIMIT 1;', [uid])
        if (data := cursor.fetchone()):
            menu += data.split(';')
        else:
            menu += [] * 14
    else:
        # Returns the selected menu
        cursor.execute('SELECT user, menu FROM menus WHERE id=?;', [mid])
        if (data := cursor.fetchone()):
            if data[0] == uid:
                menu += data[1].split(';')
            else:
                raise CAError('Menu non disponibile')
        else:
            raise CAError('Menu non trovato')

    # Fetches the prev's id
    cursor.execute('SELECT MAX(id) FROM menus WHERE user = ? AND id < ? ORDER BY id LIMIT 1;', [uid, mid])
    prev = cursor.fetchone()[0]

    # Fetches the next's id
    cursor.execute('SELECT MIN(id) FROM menus WHERE user = ? AND id > ? ORDER BY id LIMIT 1;', [uid, mid])
    next = cursor.fetchone()[0]

    return Menu(menu, prev, next)

# TODO refactor

@use_db
def update_user_menu(cursor, uid, items):
    # Checks the menu syntax
    if len(items) != 14:
        raise CAError('Menu non valido')

    # Saves the new menu
    try:
        cursor.execute('REPLACE INTO menus (user, menu) VALUES (?, ?);', [uid, ';'.join(items)])
    except MDBError:
        raise CAError('Utente sconosciuto')
