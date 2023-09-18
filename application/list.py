from . import app, cursor

from flask import render_template, request, redirect
from mariadb import Error as DBError


# Ensures the table exists
try:
    cursor.execute('SELECT * FROM list;')
except DBError:
    cursor.execute('CREATE TABLE list (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) UNIQUE NOT NULL);')


@app.route('/list', methods=['GET', 'POST'])
def list_route():
    # If the method is GET, displays the list
    if request.method == 'GET':
        cursor.execute('SELECT * FROM list;')
        items = cursor.fetchall()
        return render_template('list.html', items=items)

    # Tries to delete the selected items
    try:
        cursor.executemany('DELETE FROM list WHERE id = ?;', [(i, ) for i in request.form.keys()])
    except:
        pass

    return redirect('/list/')

@app.route('/list/add', methods=['GET', 'POST'])
def list_add_route():
    # If the method is GET, displays the form
    if request.method == 'GET':
        return render_template('add.html', date=False)

    # Otherwise, tries to add the new items
    items = [(i,) for i in request.form.values()]
    try:
        cursor.executemany('INSERT INTO list (name) VALUES (?);', items)
    except:
        pass

    return redirect('/list/')
