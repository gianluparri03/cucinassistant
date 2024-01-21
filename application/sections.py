from . import app
from .database import *
from .util import smart_route
from .account import login_required

from flask import request


@app.route('/')
@smart_route('home.html')
@login_required
def home_route(user):
    pass


@app.route('/menu/', methods=['GET', 'POST'])
@smart_route('menu.html')
@login_required
def menu_route(uid):
    # Saves the updated menu
    if request.method == 'POST':
        update_user_menu(uid, list(request.form.values()))

    # Displays the menu
    return {'menu': get_user_menu(uid)}

@app.route('/spesa/', methods=['GET', 'POST'])
@smart_route('lists.html', title='Lista della spesa')
@login_required
def shopping_route(uid):
    if request.method == 'POST':
        if 'add' in request.form:
            add_user_lists('shopping', uid, request.form['add'].split('\r\n'))
        elif 'remove' in request.form:
            remove_user_lists('shopping', uid, request.form['remove'].split('\r\n'))
        else:
            raise CAError('Richiesta sconosciuta')

    return {'list': get_user_lists('shopping', uid)}

@app.route('/idee/', methods=['GET', 'POST'])
@smart_route('lists.html', title='Lista delle idee')
@login_required
def ideas_route(uid):
    if request.method == 'POST':
        if 'add' in request.form:
            add_user_lists('ideas', uid, request.form['add'].split('\r\n'))
        elif 'remove' in request.form:
            remove_user_lists('ideas', uid, request.form['remove'].split('\r\n'))
        else:
            raise CAError('Richiesta sconosciuta')

    return {'list': get_user_lists('ideas', uid)}


@app.route('/info')
@smart_route('info.html')
def info_route():
    return {'users': get_users_number()}
