from cucinassistant.exceptions import CAError
from cucinassistant.database import db, use_user

from collections import namedtuple
from mariadb import connect, Error as MDBError


Menu = namedtuple('Menu', ['mid', 'menu', 'prev', 'next'])

@use_user
def get_menus(cursor, uid):
    cursor.execute('SELECT id FROM menus WHERE user=?', [uid])
    return [r[0] for r in cursor.fetchall()]

@use_user
def get_menu(cursor, uid, mid=None):
    if mid == None:
        # Returns the first menu
        cursor.execute('SELECT id, menu FROM menus WHERE user=? ORDER BY id DESC LIMIT 1;', [uid])
        if (data := cursor.fetchone()):
            mid = data[0]
            menu = data[1]
        else:
            mid = 0
            menu = ()
    else:
        # Returns the selected menu
        cursor.execute('SELECT menu FROM menus WHERE user=? AND id=?;', [uid, mid])
        if (data := cursor.fetchone()):
            menu = data[0]
        else:
            raise CAError('Menù non trovato')

    # Fetches the prev's id
    cursor.execute('SELECT MAX(id) FROM menus WHERE user = ? AND id < ?;', [uid, mid])
    prev = cursor.fetchone()[0]

    # Fetches the next's id
    cursor.execute('SELECT MIN(id) FROM menus WHERE user = ? AND id > ?;', [uid, mid])
    next = cursor.fetchone()[0]

    return Menu(mid, menu, prev, next)

@use_user
def create_menu(cursor, uid):
    # Creates a new menu
    prev = get_menu(uid).mid
    mid = prev + 1
    cursor.execute('INSERT INTO menus (user, id, menu, prev) VALUES (?, ?, ?, ?);', [uid, mid, ';'*13, prev or None])

    # Updates the link
    if prev > 0:
        cursor.execute('UPDATE menus SET next=? WHERE user=? AND id=?;', [mid, uid, prev])

    return mid

@use_user
def update_menu(cursor, uid, mid, menu):
    # Makes sure the menu exists (for 0 get_menu wouldn't raise anything)
    if mid == 0:
        raise CAError('Menù non trovato')
    else:
        get_menu(uid, mid)

    # Ensures the menu is valid
    if menu.count(';') != 13:
        raise CAError('Menù non valido')

    # Updates the menu
    cursor.execute('UPDATE menus SET menu=? WHERE user=? AND id=?;', [menu, uid, mid])

@use_user
def delete_menu(cursor, uid, mid):
    # Makes sure the menu exists (for 0 get_menu wouldn't raise anything)
    if mid == 0:
        raise CAError('Menù non trovato')
    else:
        menu = get_menu(uid, mid)

    # Updates the adjacent, then deletes the menu
    cursor.execute('UPDATE menus SET next=? WHERE user=? AND id=?;', [menu.next, uid, menu.prev])
    cursor.execute('UPDATE menus SET prev=? WHERE user=? AND id=?;', [menu.prev, uid, menu.next])
    cursor.execute('DELETE FROM menus WHERE user=? AND id=?;', [uid, mid])

@use_user
def duplicate_menu(cursor, uid, mid):
    # Makes sure the menu exists (for 0 get_menu wouldn't raise anything)
    if mid == 0:
        raise CAError('Menù non trovato')
    else:
        menu = get_menu(uid, mid).menu

    # Duplicates the menu
    mid2 = create_menu(uid)
    update_menu(uid, mid2, menu)
    return mid2
