from . import app
from .auth import login_required

from flask import Flask, render_template, redirect


@app.route('/impostazioni/')
@login_required
def settings_route(user):
    return render_template('settings.html')

@app.route('/impostazioni/elimina_account/', methods=['POST'])
@login_required
def delete_account_route(user):
    user.delete()
    return redirect('/login')
