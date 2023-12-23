from os import environ
from uuid import uuid4
from flask import Flask, session

# Creates the application
app = Flask(__name__)
app.secret_key = 'cucinassistant' if not environ.get('PRODUCTION') else environ['SECRET']

@app.before_request
def make_session_permanent():
    session.permanent = True

from .mail import *
from .auth import *
from .sections import *
