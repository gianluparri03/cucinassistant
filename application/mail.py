from . import config

from smtplib import SMTP
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from jinja2 import Environment, FileSystemLoader


# If it is enabled, logs in
if config['Email']['Enabled']:
    mail = SMTP(config['Email']['Server'], config['Email']['Port'])
    mail.ehlo()
    mail.starttls()
    mail.login(config['Email']['Login'], config['Email']['Password'])


# Initializes the template loader
templates = Environment(loader=FileSystemLoader("application/emails/"))


class Email:
    def __init__(self):
        self.msg = MIMEMultipart('alternative')
        self.msg['From'] = config['Email']['Address']

    def parse_template(self, filename, **data):
        if config['Email']['Enabled']:
            template = templates.get_template(filename)
            text = template.render(banner=config['Environment']['Webserver'] + '/static/banner.png', **data)
            self.msg.attach(MIMEText(text, 'html'))

    def send(self, *recipients):
        if config['Email']['Enabled']:
            for recipient in recipients:
                self.msg['To'] = recipient
                mail.sendmail(config['Email']['Address'], recipient, self.msg.as_string())


class WelcomeEmail(Email):
    def __init__(self, username, token):
        super().__init__()

        self.msg['Subject'] = 'Registrazione effettuata'
        delete_url = config['Environment']['Webserver'] + '/account/elimina/?token=' + token
        self.parse_template('welcome.html', username=username, delete_url=delete_url)

class ResetPasswordEmail(Email):
    def __init__(self, username, password):
        super().__init__()

        self.msg['Subject'] = 'Reset Password'
        self.parse_template('reset_password.html', username=username, password=password)

class ConfirmDeletionEmail(Email):
    def __init__(self, username, token):
        super().__init__()

        self.msg['Subject'] = 'Eliminazione account'
        delete_url = config['Environment']['Webserver'] + '/account/elimina/?token=' + token
        self.parse_template('confirm_deletion.html', username=username, delete_url=delete_url)

class NewVersionEmail(Email):
    def __init__(self, news):
        super().__init__()

        self.msg['Subject'] = 'Novit&agrave; della nuova versione'
        self.parse_template('new_version.html', news=news)
