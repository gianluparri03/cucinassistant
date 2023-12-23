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
            return redirect('/account/login')
    
    return inner


@app.route('/account/login', methods=['GET', 'POST'])
def signin_route(error=''):
    if request.method == 'GET' or error:
        # Returns the login page if it is a GET request (or a response to a POST)
        return render_template('login.html', error=error, can_reset=bool(config))
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or (not is_login and not data.get('email')):
            return signin_route('Dati mancanti')

        # Checks for errors
        if (err := User.check(data['username'], data['password'])):
            return signin_route(err)

        # Saves the session, then returns to the homepage
        session['username'] = data['username']
        return redirect('/')

@app.route('/account/registrazione', methods=['GET', 'POST'])
def signup_route(error=''):
    if request.method == 'GET' or error:
        # Returns the registration page if it is a GET request (or a response to a POST)
        return render_template('registrazione.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or (not is_login and not data.get('email')):
            return signup_route('Dati mancanti')

        # Checks for errors
        if (err := User.create(data['username'], data['email'], data['password'])):
            return signup_route(err)

        # Sends the welcome mail
        WelcomeEmail(data['username'], '').send(data['email'])

        # Saves the session, then returns to the homepage
        session['username'] = data['username']
        return redirect('/')

@app.route('/account/logout/')
@login_required
def logout_route(user, error=''):
    session.pop('username', None)
    return redirect('/account/login')

@app.route('/account/reset_password', methods=['GET', 'POST'])
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

@app.route('/account/')
@login_required
def settings_route(user):
    return render_template('settings.html')

@app.route('/account/elimina/', methods=['POST'])
@login_required
def delete_account_route(user):
    user.delete()
    return logout_route()

@app.route('/statistiche')
def stats_route():
    return render_template('stats.html', n_users=User.get_number())
