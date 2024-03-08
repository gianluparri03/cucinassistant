from cucinassistant.config import config

from argon2 import PasswordHasher
from mariadb import connect
from functools import wraps


def init_db(testing=False):
    global db

    # Connects to the database
    c = config['Database'] if not testing else config['DatabaseTest']
    db = connect(host=c['Host'], port=int(c['Port']), database=c['Database'], user=c['User'], password=c['Password'])
    db.autocommit = True

    # Ensures the tables have been created
    with db.cursor() as cursor:
        with open('cucinassistant/database/schema.sql') as f:
            for command in f.read().split(';')[:-1]:
                cursor.execute(command)

def use_db(func):
    @wraps(func)
    def inner(*args, **kwargs):
        global db

        # Lets the function use a db connection
        db.reconnect()
        with db.cursor() as cur:
            return func(cur, *args, **kwargs)

    return inner


# Initialize the hasher and the database connection
ph = PasswordHasher()
init_db()


from .users import *
from .menus import *
from .other import *
