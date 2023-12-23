from . import app
from .data import User
from .mail import WelcomeEmail, ResetPasswordEmail, config

from functools import wraps
from flask import render_template, request, redirect, session


def login_required(func):
    @wraps(func)
    def inner(*args, **kwargs):
        if (username := session.get('username')):
            return func(User(username), *args, **kwargs)
        else:
            return redirect('/login')
    
    return inner


@app.route('/login', methods=['GET', 'POST'])
@app.route('/registrazione', methods=['GET', 'POST'])
def login_route(error=''):
    is_login = str(request.url_rule).startswith('/login')

    if request.method == 'GET' or error:
        # Returns the login page if it is a GET request (or a response to a POST)
        return render_template('login.html', login=is_login, error=error, can_reset=bool(config))
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or (not is_login and not data.get('email')):
            return login_route('Dati mancanti')

        # Checks for errors
        if is_login:
            if (err := User.check(data['username'], data['password'])):
                return login_route(err)
        else:
            if (err := User.create(data['username'], data['email'], data['password'])):
                return login_route(err)

            # Sends the welcome mail
            WelcomeEmail(data['username'], '').send(data['email'])

        # Saves the session, then returns to the homepage
        session['username'] = data['username']
        return redirect('/')

@app.route('/reset_password', methods=['GET', 'POST'])
def reset_password_route(error=''):
    if request.method == 'GET' or error:
        # Returns the page if it is a GET request
        return render_template('reset_password.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('email'):
            return reset_password_route('Dati mancanti')
        
        # Resets the password, then sends it to the user
        if (cred := User.reset_password(data['email'])):
            ResetPasswordEmail(cred[0], cred[1]).send(data['email'])

        return reset_password_route('Email inviata. Controlla la casella di posta.')

@app.route('/logout/')
@login_required
def logout_route(user, error=''):
    session.pop('username', None)
    return redirect('/login')

@app.route('/impostazioni/')
@login_required
def settings_route(user):
    return render_template('settings.html')

@app.route('/impostazioni/elimina_account/', methods=['POST'])
@login_required
def delete_account_route(user):
    user.delete()
    return logout_route()

@app.route('/statistiche')
def stats_route():
    return render_template('stats.html', n_users=User.get_number())
