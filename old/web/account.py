@app.route('/account/elimina/', methods=['GET', 'POST'])
@smart_route('account/delete.html')
@login_required
def delete_account_route(uid):
    token = request.args.get('token')

    if request.method == 'GET':
        # Renders the page
        return {'warning': not token}
    else:
        if (data := db.get_data(uid)):
            if not token:
                # If it's the first confirm button, generates the token and sends the email
                token = db.generate_token(uid)
                delete_url = config['Environment']['Address'] + '/account/elimina/?token=' + token
                Email('Eliminazione account', 'delete_account', username=data.username, delete_url=delete_url).send(data.email)
                flash('Ti abbiamo inviato un email. Controlla la casella di posta')
            else:
                # Otherwise deletes the account
                db.delete_user(uid, token)
                Email('Eliminazione account', 'goodbye', username=data.username).send(data.email)
                flash('Account eliminato con successo')
                return logout_route()
        else:
            flash('Utente sconosciuto')

@app.route('/account/cambia_nome/', methods=['GET', 'POST'])
@smart_route('account/data_change.html', field='nome utente', field_type='text')
@login_required
def change_username_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('new'):
            flash('Dati mancanti')
            return
        
        # Tries to change the email
        db.change_username(uid, data['new'])
        flash('Nome utente cambiato con successo')

@app.route('/account/cambia_email/', methods=['GET', 'POST'])
@smart_route('account/data_change.html', field='email', field_type='email')
@login_required
def change_email_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('new'):
            flash('Dati mancanti')
            return
        
        # Tries to change the email
        db.change_email(uid, data['new'])
        flash('Email cambiata con successo')

@app.route('/account/cambia_password/', methods=['GET', 'POST'])
@smart_route('account/password_change.html')
@login_required
def change_password_route(uid):
    if request.method == 'POST':
        # Ensures the request is valid
        data = request.form
        if not data.get('old') or not data.get('new'):
            flash('Dati mancanti')
            return
        
        # Tries to change the password
        db.change_password(uid, data['old'], data['new'])
        user = db.get_data(uid)
        Email('Cambio password', 'change_password', username=user.username).send(user.email)
        flash('Password cambiata con successo')

@app.route('/account/reset_password/', methods=['GET', 'POST'])
@smart_route('account/password_reset.html')
def reset_password_route():
    if request.method == 'GET':
        return {'token': request.args.get('token')}
    else:
        # Ensures the request is valid
        data = request.form
        if not data.get('token') or not data.get('email') or not data.get('new'):
            flash('Dati mancanti')
            return
        
        # Tries to change the password
        db.reset_password(data['email'], data['token'], data['new'])
        flash('Password reimpostata con successo')
