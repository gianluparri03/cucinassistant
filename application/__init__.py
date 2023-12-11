from uuid import uuid4
from os import environ
from flask import Flask, render_template, session


# Creates the application
app = Flask(__name__)
app.secret_key = 'cucinassistant' if not environ.get('PRODUCTION') else environ['SECRET']

@app.before_request
def make_session_permanent():
    session.permanent = True

@app.after_request
def disable_cache(response):
    response.headers['Cache-Control'] = 'public, max-age=0'
    return response

@app.route('/service_worker.js')
def service_worker():
    return app.send_static_file('service_worker.js')

@app.route('/offline.html')
def offline_route():
    return render_template('offline.html')


from .auth import *
from .sections import *
