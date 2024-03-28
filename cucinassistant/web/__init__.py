from cucinassistant.config import config

from datetime import timedelta
from flask import Flask


# Creates the application
app = Flask(__name__)
app.url_map.strict_slashes = False
app.config['SECRET_KEY'] = config['Environment']['Secret']
app.config['PERMANENT_SESSION_LIFETIME'] = timedelta(days=90)
app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'
app.config['SESSION_COOKIE_SECURE'] = True


# Registers the routes
from .smart_route import *
from .account import *
from .menu import *
from .storage import *
from .lists import *
from .various import *
