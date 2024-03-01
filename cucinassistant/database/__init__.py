from cucinassistant.config import config

from argon2 import PasswordHasher
from mariadb import connect
from functools import wraps


# Connects to the database
db = connect(host=config['Database']['Hostname'], database=config['Database']['Database'], \
             user=config['Database']['Username'], password=config['Database']['Password'])
db.autocommit = True

# Ensures the tables have been created
with db.cursor() as cursor:
    with open('cucinassistant/database/schema.sql') as f:
        for command in f.read().split(';')[:-1]:
            cursor.execute(command)

# Initialize the hasher
ph = PasswordHasher()


def use_db(func):
    @wraps(func)
    def inner(*args, **kwargs):
        global db

        db.reconnect()
        with db.cursor() as cur:
            return func(cur, *args, **kwargs)

    return inner

from .users import *
from .other import *
