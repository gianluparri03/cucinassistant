from . import app
from .auth import login_required

from json import load, dump
from flask import render_template, request, redirect


@app.route('/spesa/', methods=['GET', 'POST'])
@app.route('/idee/', methods=['GET', 'POST'])
@login_required
def list_route(user, error=''):
    section = str(request.url_rule).split('/')[1]
    list = user.read_data(section)

    if request.method == 'GET' or error:
        # Displays the list if the method is GET or if there is an errror
        return render_template('list.html', list=enumerate(list), section=section, error=error)
    else:
        # Saves the updated data (if valid)
        data = request.form
        if 'add' in data.keys() and data['add']:
            list.extend(map(lambda s: s.strip(), data['add'].split('\n')))
        elif 'remove' in data.keys() and data['remove']:
            for element in set(map(lambda s: s.strip(), data['remove'].split('\n'))):
                if element and element in list:
                    list.remove(element)
                elif element:
                    return list_route('Elemento non in lista')
        else:
            return list_route('Richiesta sconosciuta')

        user.update_data(section, list)
        return redirect('.')
