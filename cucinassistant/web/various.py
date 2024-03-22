from cucinassistant.web import *
from cucinassistant import version
import cucinassistant.database as db

from flask import send_from_directory


@app.before_request
def make_session_permanent():
    session.permanent = True

@app.route('/')
@smart_route('home.html')
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
