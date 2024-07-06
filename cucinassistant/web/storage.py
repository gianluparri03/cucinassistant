from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect, url_for


@app.route('/dispensa/')
@smart_route('storage/dashboard.html')
@login_required
def storage_dashboard_route(uid):
    return {'storages': [{'id': 1, 'name': 'Frigo'}, {'id': 2, 'name': 'Pasta'}], 'str': str}

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
    if sid == 1:
        name = request.args.get('nome')
        return {'storage': db.get_storage(uid, name=name), 'filter': name, 'name': 'Frigo'}
    else:
        return {'storage': [], 'name': 'Pasta'}

@app.route('/dispensa/<int:sid>/aggiungi/')
@smart_route('storage/add.html')
@login_required
def storage_add_route_get(uid, sid):
    pass

@app.route('/dispensa/<int:sid>/aggiungi/', methods=['POST'])
@smart_route('storage/add.html')
@login_required
def storage_add_route_post(uid, sid):
    data = [a.split(';') for a in request.form.get('data', '').split('|')]
    db.append_storage(uid, data)
    return redirect('/dispensa/')

@app.route('/dispensa/<int:sid>/modifica/')
@smart_route('storage/pre_edit.html')
@login_required
def storage_pre_edit_route(uid, sid):
    name = request.args.get('nome')
    return {'storage': db.get_storage(uid, name=name), 'filter': name}

@app.route('/dispensa/<int:sid>/modifica/<int:aid>')
@smart_route('storage/edit.html')
@login_required
def storage_edit_route(uid, sid, aid):
    return {'prev': db.get_storage_article(uid, aid)}

@app.route('/dispensa/<int:sid>/modifica/<int:aid>', methods=['POST'])
@smart_route('storage/edit.html')
@login_required
def storage_edit_route_post(uid, sid, aid):
    db.edit_storage(uid, aid, [request.form.get('name'), request.form.get('expiration'), request.form.get('quantity')])
    return redirect(f'/dispensa/{sid}/')

@app.route('/dispensa/<int:sid>/rimuovi/')
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_get(uid, sid):
    name = request.args.get('nome')
    return {'storage': db.get_storage(uid, name=name), 'filter': name}


@app.route('/dispensa/<int:sid>/rimuovi/<int:aid>', methods=['POST'])
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_post(uid, sid, aid):
    db.remove_storage(uid, [aid])
    return redirect(f'/dispensa/{sid}/')

@app.route('/dispensa/<int:sid>/cerca')
@smart_route('storage/search.html')
@login_required
def storage_search_route_get(uid, sid):
    pass

@app.route('/dispensa/<int:sid>/cerca', methods=['POST'])
@smart_route('storage/search.html')
@login_required
def storage_search_route_post(uid, sid):
    return redirect(f'/dispensa/{sid}/?nome={request.form.get("name")}')
