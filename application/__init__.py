from os import environ
from mariadb import connect
from flask import Flask, render_template


# Creates the application
app = Flask(__name__)
app.url_map.strict_slashes = False

# Connects to the db
host = 'database' if environ.get('PRODUCTION') == 'true' else '127.0.0.1'
conn = connect(host=host, user='cucinassistant', password='cucinassistant', database='cucinassistant')
conn.autocommit = True
cursor = conn.cursor()


@app.route('/')
def home_route():
    return render_template('home.html')

from .list import *
from .expirations import *
