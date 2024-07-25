@use_user
def generate_token(cursor, uid):
    # Generates a new deletion token for the user
    token = token_hex(18)
    cursor.execute('UPDATE users SET token=? WHERE uid=?;', [ph.hash(token), uid])
    return token

@use_user
def delete_user(cursor, uid, token):
    # Checks if the token is valid
    cursor.execute('SELECT token FROM users WHERE uid=?;', [uid])
    try:
        if (data := cursor.fetchone())[0]:
            ph.verify(data[0], token or '')
            cursor.execute('DELETE FROM users WHERE uid=?;', [uid])
        else:
            raise VerificationError()
    except VerificationError:
        raise CAError('Errore durante la cancellazione, riprova')

@use_user
def change_username(cursor, uid, new):
    # Saves the new one
    try:
        if get_data(uid).username == new:
            return

        cursor.execute('UPDATE users SET username=? WHERE uid=?;', [new, uid])
    except MDBError:
        raise CAError('Nome utente non disponibile')

@use_user
def change_email(cursor, uid, new):
    # Saves the new one
    try:
        if get_data(uid).email == new:
            return

        cursor.execute('UPDATE users SET email=? WHERE uid=?;', [new, uid])
    except MDBError:
        raise CAError('Email non disponibile')

@use_user
def change_password(cursor, uid, old, new):
    cursor.execute('SELECT password FROM users WHERE uid=?;', [uid])
    try:
        # Check if the user is athorized, then updates it
        ph.verify(cursor.fetchone()[0], old or '')
        cursor.execute('UPDATE users SET password=? WHERE uid=?;', [ph.hash(new), uid])
    except VerificationError:
        raise CAError('Credenziali non valide')

@use_db
def reset_password(cursor, email, token, new):
    # Ensures that the user exists
    data = get_data('', email=email)
    if not data:
        raise CACritical('Utente sconosciuto')

    try:
        # Check if the user is athorized, then updates it
        if data.token:
            ph.verify(data.token, token or '')
            cursor.execute('UPDATE users SET password=?, token=NULL WHERE uid=?;', [ph.hash(new), data.uid])
        else:
            raise VerificationError()
    except VerificationError:
        raise CAError('Errore durante la reimpostazione della password')

@use_db
def get_users_emails(cursor):
    # Counts the users
    cursor.execute('SELECT email FROM users;')
    return tuple(map(lambda r: r[0], cursor.fetchall()))
