from . import app
from .mail import *
from .data import User

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


@app.route('/account/')
@login_required
def settings_route(user):
    return render_template('settings.html')

@app.route('/account/cambio_password/', methods=['GET', 'POST'])
@login_required
def change_password_route(user, error=''):
    if request.method == 'GET' or error:
        # Returns the page if it is a GET request
        return render_template('password_change.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('old') or not data.get('new'):
            return reset_password_route('Dati mancanti')
        
        # Tries to change the password
        if (err := user.change_password(data['old'], data['new'])):
            return change_password_route(err)

        return change_password_route('Password cambiata con successo.')

@app.route('/account/elimina/', methods=['GET', 'POST'])
@login_required
def delete_account_route(user, error=''):
    pre = bool(config) and not request.args.get('token')

    if request.method == 'GET' or error:
        # Renders the page
        return render_template('delete_account.html', pre=pre, error=error)
    else:
        # If it's a pre-confirm, sends an email, otherwise check
        # the token (if implemented), then deletes the user
        if pre:
            ConfirmDeletionEmail(user.username, user.get_token()).send(user.get_email())
            return delete_account_route('Email inviata. Controlla la casella di posta.')
        else:
            if bool(config) and user.get_token() != request.args.get('token'):
                return delete_account_route("Qualcosa non ha funzionato durante l'eliminazione")
            else:
                user.delete()
                return logout_route()

@app.route('/account/login', methods=['GET', 'POST'])
def signin_route(error=''):
    if request.method == 'GET' or error:
        # Returns the login page if it is a GET request (or a response to a POST)
        return render_template('signin.html', error=error, can_reset=bool(config))
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password'):
            return signin_route('Dati mancanti')

        # Checks for errors
        if (err := User.check(data['username'], data['password'])):
            return signin_route(err)

        # Saves the session, then returns to the homepage
        session['username'] = data['username']
        return redirect('/')

@app.route('/account/logout/')
@login_required
def logout_route(user, error=''):
    session.pop('username', None)
    return redirect('/account/login')

@app.route('/account/registrazione', methods=['GET', 'POST'])
def signup_route(error=''):
    if request.method == 'GET' or error:
        # Returns the registration page if it is a GET request (or a response to a POST)
        return render_template('signup.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or not data.get('email'):
            return signup_route('Dati mancanti')

        # Checks for errors
        if (err := User.create(data['username'], data['email'], data['password'])):
            return signup_route(err)

        # Sends the welcome mail
        WelcomeEmail(data['username'], User(data['username']).get_token()).send(data['email'])

        # Saves the session, then returns to the homepage
        session['username'] = data['username']
        return redirect('/')

@app.route('/account/reset_password', methods=['GET', 'POST'])
def reset_password_route(error=''):
    if request.method == 'GET' or error:
        # Returns the page if it is a GET request
        return render_template('password_reset.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('email'):
            return reset_password_route('Dati mancanti')
        
        # Resets the password, then sends it to the user
        if (cred := User.reset_password(data['email'])):
            ResetPasswordEmail(cred[0], cred[1]).send(data['email'])

        return reset_password_route('Email inviata. Controlla la casella di posta.')
