from jinja2 import Environment, FileSystemLoader
from email.mime.multipart import MIMEMultipart
from configparser import ConfigParser
from email.mime.text import MIMEText
from smtplib import SMTP


# Checks whether the mail is enabled
parser = ConfigParser()
parser.read('email.cfg')
config = parser['Email'] if 'Email' in parser else {}

# If it is, logs in
if config:
    mail = SMTP(config['Server'], config['Port'])
    mail.ehlo()
    mail.starttls()
    mail.login(config['Login'], config['Password'])


# Initializes the template loader
templates = Environment(loader=FileSystemLoader("application/emails/"))


class Email:
    def __init__(self):
        self.msg = MIMEMultipart('alternative')
        self.msg['From'] = config['Address']

    def parse_template(self, filename, **data):
        if config:
            template = templates.get_template(filename)
            text = template.render(banner=config['Webserver'] + '/static/banner.png', **data)
            self.msg.attach(MIMEText(text, 'html'))

    def send(self, *recipients):
        if config:
            for recipient in recipients:
                self.msg['To'] = recipient
                mail.sendmail(config['Address'], recipient, self.msg.as_string())


class WelcomeEmail(Email):
    def __init__(self, username, token):
        super().__init__()

        self.msg['Subject'] = 'Registrazione effettuata'
        delete_url = config['Webserver'] + '/account/elimina/?token=' + token
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
        delete_url = config['Webserver'] + '/account/elimina/?token=' + token
        self.parse_template('confirm_deletion.html', username=username, delete_url=delete_url)

class NewVersionEmail(Email):
    def __init__(self, news):
        super().__init__()

        self.msg['Subject'] = 'Novit&agrave; della nuova versione'
        self.parse_template('new_version.html', news=news)
