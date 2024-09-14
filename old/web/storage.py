from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect, url_for


@app.route('/dispensa/<int:sid>/aggiungi/')
@smart_route('storage/add.html')
@login_required
def storage_add_route_get(uid, sid):
    return {'sid': sid, 'str': str}
    pass

@app.route('/dispensa/<int:sid>/aggiungi/', methods=['POST'])
@smart_route('storage/add.html')
@login_required
def storage_add_route_post(uid, sid):
    data = [a.split(';') for a in request.form.get('data', '').split('|')]
    db.append_storage(uid, data)
    return redirect('/dispensa/')

@app.route('/dispensa/<int:sid>/<int:aid>/modifica')
@smart_route('storage/edit.html')
@login_required
def storage_edit_route(uid, sid, aid):
    return {'prev': db.get_storage_article(uid, aid), 'str': str, 'aid': aid, 'sid': sid}

@app.route('/dispensa/<int:sid>/<int:aid>/modifica', methods=['POST'])
@smart_route('storage/edit.html')
@login_required
def storage_edit_route_post(uid, sid, aid):
    db.edit_storage(uid, aid, [request.form.get('name'), request.form.get('expiration'), request.form.get('quantity')])
    return redirect(f'/dispensa/{sid}/')

@app.route('/dispensa/<int:sid>/<int:aid>/rimuovi', methods=['POST'])
@smart_route('storage/remove.html')
@login_required
def storage_remove_route_post(uid, sid, aid):
    db.remove_storage(uid, [aid])
    return redirect(f'/dispensa/{sid}/')

@app.route('/dispensa/<int:sid>/cerca')
@smart_route('storage/search.html')
@login_required
def storage_search_route_get(uid, sid):
    return {'sid': sid, 'str': str}

@app.route('/dispensa/<int:sid>/cerca', methods=['POST'])
@smart_route('storage/search.html')
@login_required
def storage_search_route_post(uid, sid):
    return redirect(f'/dispensa/{sid}/?nome={request.form.get("name")}')
