from configparser import ConfigParser
from datetime import timedelta
from flask import Flask


class CAError(Exception): pass


# Reads the config file
config = ConfigParser()
config.read('config.cfg')
if not 'Environment' in config:
    raise exit('Config file not found')

# Creates the application
app = Flask(__name__)
app.config['SECRET_KEY'] = config['Environment']['Secret']
app.config['PERMANENT_SESSION_LIFETIME'] = timedelta(days=90)
app.config['SESSION_COOKIE_SAMESITE'] = 'Lax'
app.config['SESSION_COOKIE_SECURE'] = True

# Initializes the db
from .database import init_db
init_db()

# Registers the routes
from .account import *
from .sections import *
