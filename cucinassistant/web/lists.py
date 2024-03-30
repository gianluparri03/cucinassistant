from cucinassistant.web.smart_route import smart_route
from cucinassistant.web.account import login_required
import cucinassistant.database as db
from cucinassistant.web import app

from enum import Enum
from flask import request, redirect
from werkzeug.routing import BaseConverter, ValidationError


class Sections(Enum):
    Shopping = 'spesa'
    Ideas = 'idee'

    @property
    def dbname(self):
        names = {Sections.Shopping: 'shopping', Sections.Ideas: 'ideas'}
        return names.get(self)

    @property
    def title(self):
        titles = {Sections.Shopping: 'Lista della spesa', Sections.Ideas: 'Lista delle idee'}
        return titles.get(self)

class SectionsConverter(BaseConverter):
    def to_python(self, value):
        try:
            return Sections(value)
        except ValueError:
            raise ValidationError()

    def to_url(self, obj):
        return obj.value

app.url_map.converters.update(section=SectionsConverter)


@app.route('/<section:sec>/')
@smart_route('lists/view.html')
@login_required
def lists_view_route(uid, sec):
    return {'list': db.get_list(uid, sec.dbname), 'title': sec.title, 'routename': sec.value}

@app.route('/<section:sec>/aggiungi/')
@smart_route('lists/add.html')
@login_required
def lists_add_route_get(uid, sec):
    return {'title': 'Aggiungi ' + sec.value, 'routename': sec.value}

@app.route('/<section:sec>/aggiungi/', methods=['POST'])
@smart_route('lists/add.html')
@login_required
def lists_add_route_post(uid, sec):
    data = d.split(';') if (d := request.form.get('data')) else []
    db.append_list(uid, sec.dbname, data)
    return redirect('.')

@app.route('/<section:sec>/modifica/')
@smart_route('lists/pre_edit.html')
@login_required
def shopping_pre_edit_route(uid, sec):
    return {'list': db.get_list(uid, sec.dbname), 'title': 'Modifica ' + sec.value, 'routename': sec.value}

@app.route('/<section:sec>/modifica/<int:eid>/')
@smart_route('lists/edit.html')
@login_required
def lists_edit_route_get(uid, sec, eid):
    return {'prev': db.get_list_entry(uid, sec.dbname, eid).name, 'title': 'Modifica ' + sec.value, 'routename': sec.value}

@app.route('/<section:sec>/modifica/<int:eid>/', methods=['POST'])
@smart_route('lists/edit.html')
@login_required
def lists_edit_route_post(uid, sec, eid):
    db.edit_list(uid, sec.dbname, eid, request.form.get('name'))
    return redirect('.')

@app.route('/<section:sec>/rimuovi/')
@smart_route('lists/remove.html')
@login_required
def lists_remove_route_get(uid, sec):
    return {'list': db.get_list(uid, sec.dbname), 'title': 'Rimuovi ' + sec.value, 'routename': sec.value}

@app.route('/<section:sec>/rimuovi/', methods=['POST'])
@smart_route('lists/remove.html')
@login_required
def lists_remove_route_post(uid, sec):
    db.remove_list(uid, sec.dbname, request.form.get('data').split(';'))
    return redirect('.')
