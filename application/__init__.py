from uuid import uuid4
from os import environ
from flask import Flask, render_template


# Creates the application
app = Flask(__name__)
app.secret_key = str(uuid4())

# from .list import *
# from .expirations import *
from .auth import *
from .menu import *

@app.route('/')
@login_required
def index_route(user):
    return render_template('home.html')
