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
