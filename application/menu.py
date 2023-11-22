from . import app
from .auth import login_required

from json import load, dump
from flask import render_template, request, redirect


@app.route('/menu/', methods=['GET', 'POST'])
@login_required
def menu_route(user, error=''):
    menu = user.read_data('menu')

    if request.method == 'GET' or error:
        # Displays the list if the method is GET or if there is an errror
        return render_template('menu.html', menu=menu, error=error)
    else:
        # Saves the updated menu (if valid)
        data = request.form
        if list(data.keys()) != [f'e{i}' for i in range(14)]:
            return menu_route(user, 'Menu non valido')
        else:
            user.update_data('menu', [data[f'e{i}'] for i in range(14)])
            return redirect('/menu')
