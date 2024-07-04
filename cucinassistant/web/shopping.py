from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from enum import Enum
from flask import request, redirect


@app.route('/spesa/')
@smart_route('shopping/view.html')
@login_required
def shopping_view_route(uid):
    return {'list': db.get_shopping(uid)}

@app.route('/spesa/aggiungi/')
@smart_route('shopping/add.html')
@login_required
def shopping_add_route_get(uid):
    pass

@app.route('/spesa/aggiungi/', methods=['POST'])
@smart_route('shopping/add.html')
@login_required
def shopping_add_route_post(uid):
    data = d.split(';') if (d := request.form.get('data')) else []
    db.append_shopping(uid, data)
    return redirect('/spesa/')

@app.route('/spesa/modifica/<int:eid>/')
@smart_route('shopping/edit.html')
@login_required
def shopping_edit_route_get(uid, eid):
    return {'prev': db.get_shopping_entry(uid, eid).name}

@app.route('/spesa/modifica/<int:eid>/', methods=['POST'])
@smart_route('shopping/edit.html')
@login_required
def shopping_edit_route_post(uid, eid):
    db.edit_shopping(uid, eid, request.form.get('name'))
    return redirect(f'/spesa/')

@app.route('/spesa/rimuovi/')
@smart_route('shopping/remove.html')
@login_required
def shopping_remove_route_get(uid):
    return {'list': db.get_shopping(uid)}

@app.route('/spesa/rimuovi/', methods=['POST'])
@smart_route('shopping/remove.html')
@login_required
def shopping_remove_route_post(uid):
    db.remove_shopping(uid, request.form.get('data').split(';'))
    return redirect('/spesa/')
