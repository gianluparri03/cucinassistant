from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app


@app.route('/dispensa/')
@smart_route('storage/view.html')
@login_required
def storage_view_route(uid):
    return {'storage': db.get_storage(uid)}

@app.route('/dispensa/aggiungi/')
@smart_route('storage/add.html')
@login_required
def storage_add_route(uid):
    pass

@app.route('/dispensa/modifica/')
@smart_route('storage/pre_edit.html')
@login_required
def storage_pre_edit_route(uid):
    return {'storage': data}

@app.route('/dispensa/modifica/<int:aid>')
@smart_route('storage/edit.html')
@login_required
def storage_edit_route(uid, aid):
    return {'prev': data[aid-1]}

@app.route('/dispensa/rimuovi/')
@smart_route('storage/remove.html')
@login_required
def storage_remove_route(uid):
    return {'storage': data}
