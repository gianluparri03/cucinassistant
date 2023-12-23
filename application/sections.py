from . import app
from .data import get_users_number
from .account import login_required

from json import load, dump
from datetime import datetime
from flask import render_template, request, redirect


@app.route('/')
@login_required
def home_route(user):
    return render_template('home.html')

@app.route('/spesa/', methods=['GET', 'POST'])
@app.route('/idee/', methods=['GET', 'POST'])
@app.route('/scadenze/', methods=['GET', 'POST'])
@app.route('/quantita/', methods=['GET', 'POST'])
@login_required
def list_route(user, error=''):
    section = str(request.url_rule).split('/')[1]
    list = user.read_data(section)

    if request.method == 'GET' or error:
        # Displays the list if the method is GET or if there is an errror
        return render_template('list.html', list=enumerate(list), section=section, error=error)
    else:
        # Parses the data
        data = request.form
        parse = lambda k: {s.strip() for s in data[k].split('\n') if s.strip()}

        # Adds the new elements to the old list
        if 'add' in data.keys() and data['add']:
            filter = {'scadenze': lambda d: datetime.strptime(d, '%Y-%m-%d'), 'quantita': int}
            if section in filter:
                # Ensures the items data are valid in 'scadenze' and 'quantita'
                parsed = []
                try:
                    for element in parse('add'):
                        *name, data = element.split(' ')
                        filter[section](data)
                        parsed.append((' '.join(name), data))
                except ValueError:
                    return list_route('Elemento non valido')

                list.extend(parsed)
                list.sort(key=lambda e: e[1] if section == 'scadenze' else e[0])
            else:
                list.extend(parse('add'))
                if section != 'menu':
                    list.sort()

        # Or remove some items
        elif 'remove' in data.keys() and data['remove'] and section != 'quantita':
            # Ensures the ids are integers
            elements = parse('remove')
            if not all(e.isnumeric() for e in elements):
                return list_route('Elemento non valido')

            # Ensures the ids are valid
            elements = sorted(map(int, elements), reverse=True)
            if elements[0] >= len(list):
                return list_route('Elemento non in lista')

            # Removes them
            for element in elements:
                list.pop(element)

        # Or remove some items
        elif all((k in data.keys() and data[k]) for k in ('edit-id', 'edit-delta')) and section == 'quantita':
            index = data['edit-id'].strip()
            delta = data['edit-delta'].strip()

            # Ensures the values are valid
            if not index.isnumeric():
                return list_route('Elemento non valido')
            elif int(index) >= len(list):
                return list_route('Elemento non in lista')
            elif not delta[1:].isnumeric() or (delta[0] != '-' and delta[0] != '+'):
                return list_route('Valore non valido')

            # Updates the item
            index = int(index)
            if (new := eval(list[index][1] + delta)) > 0:
                list[index][1] = str(new)
            else:
                list.pop(index)

        else:
            return list_route('Richiesta sconosciuta')

        user.update_data(section, list)
        return redirect('.')

@app.route('/menu/', methods=['GET', 'POST'])
@login_required
def menu_route(user, error=''):
    menu = user.read_data('menu')

    if request.method == 'GET' or error:
        # Displays the menu if the method is GET or if there is an errror
        return render_template('menu.html', menu=menu, error=error)
    else:
        # Saves the updated menu (if valid)
        data = request.form
        if list(data.keys()) != [f'e{i}' for i in range(14)]:
            return menu_route(user, 'Menu non valido')
        else:
            user.update_data('menu', [data[f'e{i}'] for i in range(14)])
            return redirect('.')

@app.route('/statistiche')
def stats_route():
    return render_template('stats.html', n_users=get_users_number())
