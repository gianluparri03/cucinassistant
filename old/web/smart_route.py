from cucinassistant.exceptions import CAError, CACritical
from cucinassistant.config import config 

from flask import render_template, Response, flash, redirect, request
from functools import wraps


# Inside a smart route, any CAError (and CACritical) is automatically catched and
# rendered inside the specified template. If all goes well, the return
# value of the function is rendered: if it's a response, that response
# will be sent to the client; if it's a dict, it will be used in the
# template, and finally, if it's a NoneType, the template will be rendered
# on its own.
def smart_route(template, **data):
    def inner(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            try:
                res = func(*args, **kwargs)
                match type(res).__name__:
                    case 'Response':
                        return res

                    case 'dict':
                        data.update(res)

            except CACritical as err:
                flash(str(err))
                return redirect('/account/esci')

            except CAError as err:
                flash(str(err))
                return redirect('/')

            return render_template(template,
                                   config=config,
                                   is_hx=lambda: request.headers.get('HX-Request', False),
                                   **data)
        return wrapper
    return inner
