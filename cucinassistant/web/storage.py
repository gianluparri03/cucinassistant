from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect, url_for


@app.route('/dispensa/')
@smart_route('storage/dashboard.html')
@login_required
def storage_dashboard_route(uid):
    return {'storages': [{'id': 1, 'name': 'Frigo'}], 'str': str}

@app.route('/dispensa/crea/')
@smart_route('storage/new.html')
@login_required
def storage_new_route_get(uid):
    pass

@app.route('/dispensa/crea/', methods=["POST"])
@login_required
def storage_new_route_post(uid):
    return redirect('/dispensa/')

@app.route('/dispensa/<int:sid>/')
@smart_route('storage/view.html')
@login_required
def storage_view_route(uid, sid):
    name = request.args.get('nome')
    return {'storage': db.get_storage(uid, name=name), 'filter': name, 'name': 'Frigo'}

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
    return redirect('/dispensa/')

@app.route('/dispensa/modifica/')
@smart_route('storage/pre_edit.html')
@login_required
def storage_pre_edit_route(uid):
    name = request.args.get('nome')
    return {'storage': db.get_storage(uid, name=name), 'filter': name}

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
    return redirect('/dispensa/modifica/')

@app.route('/dispensa/rimuovi/')
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_get(uid):
    name = request.args.get('nome')
    return {'storage': db.get_storage(uid, name=name), 'filter': name}


@app.route('/dispensa/rimuovi/', methods=['POST'])
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_post(uid):
    db.remove_storage(uid, request.form.get('data').split(';'))
    return redirect('.')

@app.route('/dispensa/cerca')
@smart_route('storage/search.html')
@login_required
def storage_search_route_get(uid):
    return {'page': request.args.get('pagina', '')}

@app.route('/dispensa/cerca', methods=['POST'])
@smart_route('storage/search.html')
@login_required
def storage_search_route_post(uid):
    pages = {'': 'storage_view_route', 'rimuovi': 'storage_remove_route_get', 'modifica': 'storage_pre_edit_route'}
    page = pages.get(request.form.get('page'), pages[''])
    return redirect(url_for(page, nome=request.form.get('name', '')))
