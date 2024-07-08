from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from flask import request, redirect


@app.route('/menu/')
@smart_route('menu/dashboard.html')
@login_required
def menu_dashboard_route(uid):
    return {'menus': db.get_menus(uid), 'str': str}

@app.route('/menu/<int:mid>/')
@smart_route('menu/view.html')
@login_required
def menu_view_route(uid, mid=None):
    return {'menu': db.get_menu(uid, mid), 'str':str}

@app.route('/menu/<int:mid>/modifica')
@smart_route('menu/edit.html')
@login_required
def menu_edit_route_get(uid, mid):
    return {'menu': db.get_menu(uid, mid), 'str': str}

@app.route('/menu/<int:mid>/modifica', methods=['POST'])
@smart_route('menu/edit.html')
@login_required
def menu_edit_route_post(uid, mid):
    data = ';'.join(request.form.get(f'entry-{i}', '') for i in range(14))
    db.update_menu(uid, mid, data)
    return redirect(f'/menu/{mid}')

@app.route('/menu/nuovo', methods=['POST'])
@smart_route('menu/view.html')
@login_required
def menu_create_route(uid):
    mid = db.create_menu(uid)
    return redirect(f'/menu/{mid}/modifica')

@app.route('/menu/<int:mid>/elimina', methods=['POST'])
@smart_route('menu/view.html')
@login_required
def menu_delete_route(uid, mid):
    db.delete_menu(uid, mid)
    return redirect('/menu')

@app.route('/menu/<int:mid>/clona', methods=['POST'])
@smart_route('menu/view.html')
@login_required
def menu_clone_route(uid, mid):
    mid = db.duplicate_menu(uid, mid)
    return redirect(f'/menu/{mid}/modifica')
