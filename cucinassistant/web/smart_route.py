from cucinassistant.config import config 
from cucinassistant.exceptions import CAError

from flask import render_template, Response
from functools import wraps


# Inside a smart route, any CAError is automatically catched and
# rendered inside the specified template. If all goes well, the return
# value of the function is rendered: if it's a response, that response
# will be sent to the client; if it's a string, that message will be shown;
# if it's a dict, it will be used in the template, and finally, if it's a
# NoneType, the template will be rendered on its own.
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
                        show = ''
                        data.update(res)

                    case 'str':
                        show = res

                    case 'NoneType':
                        show = ''
            except CAError as err:
                show = str(err)

            return render_template(template, **data, show=show, config=config)
        return wrapper
    return inner
