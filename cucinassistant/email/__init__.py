from cucinassistant.config import config
from cucinassistant.database import get_users_emails

from smtplib import SMTP
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from jinja2 import Environment, FileSystemLoader


templates = Environment(loader=FileSystemLoader('cucinassistant/email/templates'))


class Email:
    def __init__(self, subject, template_name, **data):
        if config['Email']['Enabled']:
            # Initializes the email
            self.msg = MIMEMultipart('alternative')
            self.msg['From'] = config['Email']['Address']
            self.msg['Subject'] = subject

            # Parses the template
            template = templates.get_template(template_name + '.html')
            text = template.render(banner=config['Environment']['Address'] + '/static/img/banner.png', **data)
            self.msg.attach(MIMEText(text, 'html'))

    def connect(self):
        # Connects to the server
        conn = SMTP(config['Email']['Server'], config['Email']['Port'])
        conn.ehlo()
        conn.starttls()
        conn.login(config['Email']['Login'], config['Email']['Password'])
        return conn

    def send(self, recipient):
        if config['Email']['Enabled']:
            # Sends the email to the recipient
            with self.connect() as conn:
                self.msg['To'] = recipient
                conn.sendmail(config['Email']['Address'], recipient, self.msg.as_string())

    def broadcast(self):
        if config['Email']['Enabled']:
            # Sends the email to everyone
            with self.connect() as conn:
                recipients = get_users_emails()
                self.msg['Bcc'] = ', '.join(recipients)
                conn.sendmail(config['Email']['Address'], recipients, self.msg.as_string())
                return len(recipients)
