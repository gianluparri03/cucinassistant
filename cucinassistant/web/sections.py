from . import app, Version
from .database import *
from .util import smart_route
from .account import login_required, is_logged

from flask import request, redirect, send_from_directory



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
