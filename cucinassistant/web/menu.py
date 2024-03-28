from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect


@app.route('/menu/')
@app.route('/menu/<int:mid>/')
@smart_route('menu/view.html')
@login_required
def menu_view_route(uid, mid=None):
    return {'menu': db.get_menu(uid, mid)}

@app.route('/menu/<int:mid>/modifica')
@smart_route('menu/edit.html')
@login_required
def menu_edit_route_get(uid, mid):
    return {'menu': db.get_menu(uid, mid)}

@app.route('/menu/<int:mid>/modifica', methods=['POST'])
@smart_route('menu/edit.html')
@login_required
def menu_edit_route_post(uid, mid):
    db.update_menu(uid, mid, request.form.get('data', ''))
    return redirect('.')

@app.route('/menu/crea')
@login_required
def menu_create_route(uid):
    mid = db.create_menu(uid)
    flash('Men&ugrave; creato correttamente.')
    return redirect(f'/menu/{mid}')

@app.route('/menu/<int:mid>/elimina', methods=['POST'])
@login_required
def menu_delete_route(uid, mid):
    db.delete_menu(uid, mid)
    flash('Men&ugrave; eliminato correttamente.')
    return redirect('/menu')

@app.route('/menu/<int:mid>/clona')
@login_required
def menu_clone_route(uid, mid):
    mid = db.duplicate_menu(uid, mid)
    flash('Men&ugrave; clonato correttamente')
    return redirect(f'/menu/{mid}')
