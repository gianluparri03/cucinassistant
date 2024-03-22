from cucinassistant.exceptions import CAError
import cucinassistant.database as db
from cucinassistant.web import *

from flask import request, redirect


def lists_view_route(uid, db_name):
    return {'list': db.get_list(uid, db_name)}

@app.route('/spesa/')
@smart_route('lists/view.html', title='Lista della spesa')
@login_required
def shopping_view_route(uid):
    return lists_view_route(uid, 'shopping')

@app.route('/idee/')
@smart_route('lists/view.html', title='Idee')
@login_required
def ideas_view_route(uid):
    return lists_view_route(uid, 'ideas')


def lists_add_route_get(uid, db_name):
    return {'list': db.get_list(uid, db_name)}

@app.route('/spesa/aggiungi/')
@smart_route('lists/add.html', title='Aggiungi spesa')
@login_required
def shopping_add_route_get(uid):
    return lists_add_route_get(uid, 'shopping')

@app.route('/idee/aggiungi/')
@smart_route('lists/add.html', title='Aggiungi idee')
@login_required
def ideas_add_route_get(uid):
    return lists_add_route_get(uid, 'ideas')


def lists_add_route_post(uid, db_name):
    data = d.split(';') if (d := request.form.get('data')) else []
    db.append_list(uid, db_name, data)
    return redirect('.')

@app.route('/spesa/aggiungi/', methods=['POST'])
@smart_route('lists/add.html', title='Aggiungi spesa')
@login_required
def shopping_add_route_post(uid):
    return lists_add_route_post(uid, 'shopping')

@app.route('/idee/aggiungi/', methods=['POST'])
@smart_route('lists/add.html', title='Aggiungi idee')
@login_required
def ideas_add_route_post(uid):
    return lists_add_route_post(uid, 'ideas')


def lists_edit_route_get(uid, db_name, eid):
    return {'prev': db.get_list_entry(uid, db_name, eid)}

@app.route('/spesa/modifica/<int:eid>/')
@smart_route('lists/edit.html', title='Modifica spesa')
@login_required
def shopping_edit_route_get(uid, eid):
    return lists_edit_route_get(uid, 'shopping', eid)

@app.route('/idee/modifica/<int:eid>/')
@smart_route('lists/edit.html', title='Modifica idee')
@login_required
def ideas_edit_route_get(uid, eid):
    return lists_edit_route_get(uid, 'ideas', eid)


def lists_edit_route_post(uid, db_name, eid):
    db.edit_list(uid, db_name, eid, request.form.get('name'))
    return redirect('..')

@app.route('/spesa/modifica/<int:eid>/', methods=['POST'])
@smart_route('lists/edit.html', title='Modifica spesa')
@login_required
def shopping_edit_route_post(uid, eid):
    return lists_edit_route_post(uid, 'shopping', eid)

@app.route('/idee/modifica/<int:eid>/', methods=['POST'])
@smart_route('lists/edit.html', title='Modifica idee')
@login_required
def ideas_edit_route_post(uid, eid):
    return lists_edit_route_post(uid, 'ideas', eid)


@app.route('/spesa/rimuovi/')
@smart_route('lists/remove.html', title='Rimuovi spesa')
@login_required
def shopping_remove_route_get(uid):
    return lists_view_route(uid, 'shopping')

@app.route('/idee/rimuovi/')
@smart_route('lists/remove.html', title='Rimuovi idee')
@login_required
def ideas_remove_route_get(uid):
    return lists_view_route(uid, 'ideas')


def lists_remove_route_post(uid, db_name):
    db.remove_list(uid, db_name, request.form.get('data').split(';'))
    return redirect('.')

@app.route('/spesa/rimuovi/', methods=['POST'])
@smart_route('lists/remove.html', title='Modifica spesa')
@login_required
def shopping_remove_route_post(uid):
    return lists_remove_route_post(uid, 'shopping')

@app.route('/idee/rimuovi/', methods=['POST'])
@smart_route('lists/remove.html', title='Modifica idee')
@login_required
def ideas_remove_route_post(uid):
    return lists_remove_route_post(uid, 'ideas')
