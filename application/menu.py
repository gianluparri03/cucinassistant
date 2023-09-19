from . import app, cursor

from flask import render_template, request, redirect
from mariadb import Error as DBError


# Ensures the table exists
try:
    cursor.execute('SELECT * FROM menu;')
except DBError:
    cursor.execute('CREATE TABLE menu (id INT PRIMARY KEY, name VARCHAR(100));')


@app.route('/menu', methods=['GET', 'POST'])
def menu_route():
    def get_items():
        cursor.execute('SELECT * FROM menu;')
        return dict(cursor.fetchall())

    # If the method is GET, displays the list
    if request.method == 'GET':
        return render_template('menu.html', items=get_items())

    # Tries to save the new menu
    try:
        cursor.execute('DELETE FROM menu;')
        cursor.executemany('INSERT INTO menu (id, name) VALUES (?, ?);', list(request.form.items()))
    except Exception as e:
        return render_template('menu.html', items=get_items(), error=str(e))

    return redirect('/menu/')
