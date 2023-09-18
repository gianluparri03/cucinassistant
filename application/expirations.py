from . import app, cursor

from flask import render_template, request, redirect
from mariadb import Error as DBError


# Ensures the table exists
try:
    cursor.execute('SELECT * FROM expirations;')
except DBError:
    cursor.execute('CREATE TABLE expirations (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) UNIQUE NOT NULL, date DATE NOT NULL);')


@app.route('/expirations', methods=['GET', 'POST'])
def expirations_route():
    # If the method is GET, displays the list
    if request.method == 'GET':
        cursor.execute('SELECT * FROM expirations ORDER BY date;')
        items = [(i[0], i[1], i[2].strftime('%d/%b/%y')) for i in cursor.fetchall()]
        return render_template('expirations.html', items=items)

    # Tries to delete the selected items
    try:
        cursor.executemany('DELETE FROM expirations WHERE id = ?;', [(i, ) for i in request.form.keys()])
    except:
        pass

    return redirect('/expirations/')

@app.route('/expirations/add', methods=['GET', 'POST'])
def expirations_add_route():
    # If the method is GET, displays the form
    if request.method == 'GET':
        return render_template('add.html', date=True)

    # Otherwise, tries to add the new items
    try:
        items = [(request.form[i], request.form[i.replace('-name', '-date')]) for i in request.form if i.endswith('-name')]
        cursor.executemany('INSERT INTO expirations (name, date) VALUES (?, ?);', items)
    except:
        pass

    return redirect('/expirations/')
