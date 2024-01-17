from os import environ
from uuid import uuid4
from sqlite3 import connect
from datetime import timedelta
from flask import Flask, session


# Creates the application
app = Flask(__name__)
app.secret_key = 'cucinassistant' if not environ.get('PRODUCTION') else environ['SECRET']
app.config['PERMANENT_SESSION_LIFETIME'] = timedelta(days=90)
app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'
app.config['SESSION_COOKIE_SECURE'] = True

# Connects to the database
db = connect('cucinassistant.db', check_same_thread=False, isolation_level=None)

@app.before_request
def make_session_permanent():
    session.permanent = True

    # Creates the users table
    cursor = db.cursor()
    cursor.execute('''CREATE TABLE IF NOT EXISTS users (
                    username TEXT NOT NULL PRIMARY KEY,
                    password TEXT NOT NULL,
                    email TEXT NOT NULL UNIQUE,
                    token TEXT,
                    newsletter BOOLEAN DEFAULT TRUE);''')


from .mail import *
from .account import *
from .sections import *
