from cucinassistant.exceptions import CAError
from cucinassistant.email import Email
import cucinassistant.database as db
from cucinassistant.web import *

from functools import wraps
from flask import request, redirect, session



def login_required(func):
    @wraps(func)
    def inner(*args, **kwargs):
        if (uid := session.get('user')):
            return func(uid, *args, **kwargs)
        else:
            return redirect('/account/accedi')
    
    return inner

def is_logged():
    return 'user' in session


@app.route('/account/impostazioni')
@smart_route('account/settings.html')
@login_required
def settings_route(uid):
    pass


@app.route('/account/accedi', methods=['GET', 'POST'])
@smart_route('account/signin.html')
def signin_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password'):
            raise CAError('Dati mancanti')

        # Signs in the user
        uid = db.login_user(data['username'], data['password'])

        # Saves the session, then returns to the homepage
        session['user'] = uid
        return redirect('/')

@app.route('/account/registrati', methods=['GET', 'POST'])
@smart_route('account/signup.html')
def signup_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('username') or not data.get('password') or not data.get('email'):
            raise CAError('Dati mancanti')

        # Signs up the user and sends it the welcome email
        uid = db.create_user(data['username'], data['email'], data['password'])
        Email('Registrazione effettuata', 'welcome', username=data['username']).send(data['email'])

        # Saves the session, then returns to the homepage
        session['user'] = uid
        return redirect('/')

@app.route('/account/esci/')
@login_required
def logout_route(uid):
    session.pop('user', None)
    return redirect('/account/accedi')

@app.route('/account/elimina/', methods=['GET', 'POST'])
@smart_route('account/delete.html')
@login_required
def delete_account_route(uid):
    token = request.args.get('token')

    if request.method == 'GET':
        # Renders the page
        return {'warning': not token}
    else:
        if (data := get_user_data(uid)):
            if not token:
                # If it's the first confirm button, generates the token and sends the email
                token = db.generate_user_token(uid)
                delete_url = config['Environment']['Address'] + '/account/elimina/?token=' + token
                Email('Eliminazione account', 'delete_account', username=data['username'], delete_url=delete_url).send(data['email'])
                return 'Ti abbiamo inviato un email. Controlla la casella di posta'
            else:
                # Otherwise deletes the account
                db.delete_user(uid, token)
                Email('Eliminazione account', 'goodbye', username=data['username']).send(data['email'])
                return logout_route()
        else:
            return 'Utente sconosciuto'

@app.route('/account/cambia_nome/', methods=['GET', 'POST'])
@smart_route('account/data_change.html', field='nome', field_type='text')
@login_required
def change_username_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('new'):
            return 'Dati mancanti'
        
        # Tries to change the email
        db.change_user_username(uid, data['new'])
        return 'Nome utente cambiato con successo'

@app.route('/account/cambia_email/', methods=['GET', 'POST'])
@smart_route('account/data_change.html', field='email', field_type='email')
@login_required
def change_email_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('new'):
            return 'Dati mancanti'
        
        # Tries to change the email
        db.change_user_email(uid, data['new'])
        return 'Email cambiata con successo'

@app.route('/account/cambia_password/', methods=['GET', 'POST'])
@smart_route('account/password_change.html', type='change')
@login_required
def change_password_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('old') or not data.get('new'):
            return 'Dati mancanti'
        
        # Tries to change the password
        db.change_user_password(uid, data['old'], data['new'])
        user = db.get_user_data(uid)
        Email('Cambio password', 'change_password', username=user['username']).send(user['email'])
        return 'Password cambiata con successo'

@app.route('/account/recupera_password', methods=['GET', 'POST'])
@smart_route('account/password_recover.html')
def recover_password_route():
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('email'):
            return 'Dati mancanti'
        
        # Sends the user a token to reset the password
        if (data := db.get_user_data('', email=data['email'])):
            token = db.generate_user_token(data['uid'])
            change_url = config['Environment']['Address'] + '/account/reset_password/?token=' + token
            Email('Recupera password', 'reset_password', username=data['username'], change_url=change_url).send(data['email'])

        return 'Ti abbiamo inviato un email. Controlla la casella di posta'

@app.route('/account/reset_password/', methods=['GET', 'POST'])
@smart_route('account/password_change.html', type='reset')
def reset_password_route():
    if request.method == 'GET':
        return {'token': request.args.get('token')}
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('token') or not data.get('email') or not data.get('new'):
            return 'Dati mancanti'
        
        # Tries to change the password
        db.reset_user_password(data['email'], data['token'], data['new'])
        return 'Password reimpostata con successo'
