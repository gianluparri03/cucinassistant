from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect


@app.route('/dispensa/')
@smart_route('storage/view.html')
@login_required
def storage_view_route(uid):
    return {'storage': db.get_storage(uid)}

@app.route('/dispensa/aggiungi/')
@smart_route('storage/add.html')
@login_required
def storage_add_route_get(uid):
    pass

@app.route('/dispensa/aggiungi/', methods=['POST'])
@smart_route('storage/add.html')
@login_required
def storage_add_route_post(uid):
    data = [a.split(';') for a in request.form.get('data', '').split('|')]
    db.append_storage(uid, data)
    return redirect('.')

@app.route('/dispensa/modifica/')
@smart_route('storage/pre_edit.html')
@login_required
def storage_pre_edit_route(uid):
    return {'storage': db.get_storage(uid)}

@app.route('/dispensa/modifica/<int:aid>')
@smart_route('storage/edit.html')
@login_required
def storage_edit_route(uid, aid):
    return {'prev': db.get_storage_article(uid, aid)}

@app.route('/dispensa/modifica/<int:aid>', methods=['POST'])
@smart_route('storage/edit.html')
@login_required
def storage_edit_route_post(uid, aid):
    db.edit_storage(uid, aid, [request.form.get('name'), request.form.get('expiration'), request.form.get('quantity')])
    return redirect('.')

@app.route('/dispensa/rimuovi/')
@smart_route('storage/remove.html')
@login_required
def storage_remove_route(uid):
    return {'storage': db.get_storage(uid)}

@app.route('/dispensa/rimuovi/', methods=['POST'])
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_post(uid):
    db.remove_storage(uid, request.form.get('data').split(';'))
    return redirect('.')
