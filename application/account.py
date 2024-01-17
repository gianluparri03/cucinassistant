from . import app
from .mail import *
from .database import *

from functools import wraps
from flask import render_template, request, redirect, session


def login_required(func):
    @wraps(func)
    def inner(*args, **kwargs):
        if (username := session.get('username')):
            return func(username, *args, **kwargs)
        else:
            return redirect('/account/login')
    
    return inner


@app.route('/account/')
@login_required
def account_route(user):
    return render_template('account.html')


@app.route('/account/login', methods=['GET', 'POST'])
def signin_route(error=''):
    if request.method == 'GET' or error:
        # Returns the login page if it is a GET request (or a response to a POST)
        return render_template('signin.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password'):
            return signin_route('Dati mancanti')

        try:
            # Signs in the user
            login_user(data['username'], data['password'])

            # Saves the session, then returns to the homepage
            session['username'] = data['username']
            return redirect('/')

        # Displays the error if it occurs
        except CAError as e:
            return signin_route(e)

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

        try:
            # Signs up the user
            token = create_user(data['username'], data['email'], data['password'])

            # Sends the welcome mail
            WelcomeEmail(data['username'], token).send(data['email'])

            # Saves the session, then returns to the homepage
            session['username'] = data['username']
            return redirect('/')

        # Displays the error if it occurs
        except CAError as e:
            return signup_route(e)

@app.route('/account/logout/')
@login_required
def logout_route(username, error=''):
    session.pop('username', None)
    return redirect('/account/login')


@app.route('/account/elimina/', methods=['GET', 'POST'])
@login_required
def delete_account_route(username, error=''):
    token = request.args.get('token')

    if request.method == 'GET' or error:
        # Renders the page
        return render_template('delete_account.html', warning=not token, error=error)
    else:
        if not token:
            # If it's the first confirm button, generates the token, then
            # sends the email
            email = get_user_email(username)
            token = generate_user_token(username)
            ConfirmDeletionEmail(username, token).send(email)
            return delete_account_route('Ti abbiamo inviato un email. Controlla la casella di posta.')
        else:
            # Otherwise deletes the account
            try:
                delete_user(username, token)
                return logout_route()
            except CAError as err:
                return delete_account_route(str(err))


@app.route('/account/cambio_password/', methods=['GET', 'POST'])
@login_required
def change_password_route(username, error=''):
    if request.method == 'GET' or error:
        # Returns the page if it is a GET request
        return render_template('password_change.html', error=error)
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('old') or not data.get('new'):
            return reset_password_route('Dati mancanti')
        
        # Tries to change the password
        try:
            change_user_password(username, data['old'], data['new'])
            return change_password_route('Password cambiata con successo.')
        except CAError as err:
            return change_password_route(str(err))

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
        if (cred := reset_user_password(data['email'])):
            ResetPasswordEmail(cred[0], cred[1]).send(data['email'])

        return reset_password_route('Ti abbiamo inviato un email. Controlla la casella di posta.')
