from cucinassistant.web.account import login_required, is_logged
from cucinassistant.web.smart_route import smart_route
from cucinassistant.version import version
from cucinassistant.config import config
import cucinassistant.database as db
from cucinassistant.web import app

from flask import send_from_directory, session, redirect


@app.before_request
def make_session_permanent():
    session.permanent = True

@app.route('/')
@smart_route('other/home.html', get_username=lambda: db.get_data(session['user']).username)
@login_required
def home_route(uid):
    pass

@app.route('/info')
@smart_route('other/info.html', get_users_no=db.get_users_number, version=version, is_logged=is_logged)
def info_route():
    pass

@app.route('/privacy')
@smart_route('other/privacy.html', is_logged=is_logged)
def privacy_route():
    pass

@app.route('/favicon.ico')
def favicon_route():
    return send_from_directory('static', 'img/logo.png')

@app.route('/guida')
def tutorial_route():
    return redirect(config['Various']['Tutorial'])

@app.route('/serviceworker.js')
def serviceworker_route():
    return send_from_directory('static', 'js/serviceworker.js', mimetype='text/javascript')
