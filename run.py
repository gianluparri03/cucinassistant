from application import app
from os import environ


debug = not environ.get('PRODUCTION')
app.run(host='0.0.0.0', port=8080, debug=debug)
