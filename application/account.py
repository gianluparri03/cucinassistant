from . import app, CAError
from .util import smart_route
from .database import *
from .mail import *

from functools import wraps
from flask import request, redirect, session


@app.before_request
def make_session_permanent():
    session.permanent = True

def login_required(func):
    @wraps(func)
    def inner(*args, **kwargs):
        if (uid := session.get('user')):
            return func(uid, *args, **kwargs)
        else:
            return redirect('/account/accedi')
    
    return inner


@app.route('/account/')
@smart_route('account.html')
@login_required
def account_route(uid):
    pass


@app.route('/account/accedi', methods=['GET', 'POST'])
@smart_route('signin.html')
def signin_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password'):
            raise CAError('Dati mancanti')

        # Signs in the user
        uid = login_user(data['username'], data['password'])

        # Saves the session, then returns to the homepage
        session['user'] = uid
        return redirect('/')

@app.route('/account/registrazione', methods=['GET', 'POST'])
@smart_route('signup.html')
def signup_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or not data.get('email'):
            raise CAError('Dati mancanti')

        # Signs up the user
        uid = create_user(data['username'], data['email'], data['password'])

        # Sends the welcome mail
        WelcomeEmail(data['username']).send(data['email'])

        # Saves the session, then returns to the homepage
        session['user'] = uid
        return redirect('/')

@app.route('/account/logout/')
@login_required
def logout_route(uid):
    session.pop('user', None)
    return redirect('/account/accedi')


@app.route('/account/elimina/', methods=['GET', 'POST'])
@smart_route('delete_account.html')
@login_required
def delete_account_route(uid):
    token = request.args.get('token')

    if request.method == 'GET':
        # Renders the page
        return {'token': not token}

    else:
        if not token:
            # If it's the first confirm button, generates the token, then
            # sends the email
            email = get_user_email(uid)
            username = get_user_username(uid)
            token = generate_user_token(uid)
            ConfirmDeletionEmail(username, token).send(email)
            return 'Ti abbiamo inviato un email. Controlla la casella di posta.'
        else:
            # Otherwise deletes the account
            delete_user(uid, token)
            return logout_route()


@app.route('/account/cambio_password/', methods=['GET', 'POST'])
@smart_route('password_change.html')
@login_required
def change_password_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('old') or not data.get('new'):
            raise CAError('Dati mancanti')
        
        # Tries to change the password
        change_user_password(uid, data['old'], data['new'])
        return 'Password cambiata con successo.'

@app.route('/account/reset_password', methods=['GET', 'POST'])
@smart_route('password_change.html')
def reset_password_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('email'):
            raise CAError('Dati mancanti')
        
        # Resets the password, then sends it to the user
        if (cred := reset_user_password(data['email'])):
            ResetPasswordEmail(cred[0], cred[1]).send(data['email'])

        return 'Ti abbiamo inviato un email. Controlla la casella di posta.'
