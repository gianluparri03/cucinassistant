from . import app
from .auth import login_required
from .user import User

from flask import render_template


@app.route('/')
@login_required
def home_route(user):
    return render_template('home.html')

@app.route('/privacy')
def privacy_route():
    return render_template('privacy.html')

@app.route('/statistiche')
def stats_route():
    return render_template('stats.html', n_users=User.get_number())
