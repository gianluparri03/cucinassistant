from configparser import ConfigParser


config = ConfigParser()
config.read('config.cfg')
if not 'Environment' in config:
    raise exit('Config file not found')
