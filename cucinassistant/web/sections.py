from . import app, Version
from .database import *
from .util import smart_route
from .account import login_required, is_logged

from flask import request, redirect, send_from_directory


@app.route('/.well-known/assetlinks.json')
def assets_route():
    return send_from_directory('static', 'assets.json')


# TODO split into multiple files


@app.route('/menu/', methods=['GET', 'POST'])
@smart_route('menu.html', get_menu=login_required(get_user_menu))
@login_required
def menu_route(uid):
    if request.method == 'POST':
        update_user_menu(uid, list(request.form.values()))
        return redirect('.')

    return {'edit': 'modifica' in request.args}

@app.route('/dispensa/', methods=['GET', 'POST'])
@smart_route('storage.html', get_storage=login_required(get_user_storage))
@login_required
def storage_route(uid):
    if request.method == 'POST':
        if 'add' in request.form:
            add_user_storage(uid, [s.split(';') for s in request.form['add'].split('\r\n')])
        elif 'edit' in request.form:
            data = request.form['edit'].split(';') + ['+0']
            edit_user_storage(uid, data[0], data[1])
            return redirect('?modifica')
        elif 'remove' in request.form:
            remove_user_storage(uid, request.form['remove'].split('\r\n'))
        elif request.form:
            raise CAError('Richiesta sconosciuta')

        return redirect('.')

    return {'add': 'aggiungi' in request.args, 'edit': 'modifica' in request.args, 'del': 'rimuovi' in request.args}

@app.route('/spesa/', methods=['GET', 'POST'])
@smart_route('lists.html', title='Lista della spesa', get_list=lambda: login_required(get_user_lists)('shopping'))
@login_required
def shopping_route(uid):
    if request.method == 'POST':
        if 'add' in request.form:
            add_user_lists(uid, 'shopping', request.form['add'].split('\r\n'))
        elif 'remove' in request.form:
            remove_user_lists(uid, 'shopping', request.form['remove'].split('\r\n'))
        elif request.form:
            raise CAError('Richiesta sconosciuta')

        return redirect('.')

    return {'add': 'aggiungi' in request.args, 'del': 'rimuovi' in request.args}

@app.route('/idee/', methods=['GET', 'POST'])
@smart_route('lists.html', title='Lista delle idee', get_list=lambda: login_required(get_user_lists)('ideas'))
@login_required
def ideas_route(uid):
    if request.method == 'POST':
        if 'add' in request.form:
            add_user_lists(uid, 'ideas', request.form['add'].split('\r\n'))
        elif 'remove' in request.form:
            remove_user_lists(uid, 'ideas', request.form['remove'].split('\r\n'))
        elif request.form:
            raise CAError('Richiesta sconosciuta')

        return redirect('.')

    return {'add': 'aggiungi' in request.args, 'del': 'rimuovi' in request.args}
