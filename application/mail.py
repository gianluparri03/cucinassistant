from . import config
from .database import get_users_emails

from smtplib import SMTP
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from jinja2 import Environment, FileSystemLoader


# Initializes the template loader
templates = Environment(loader=FileSystemLoader("application/emails/"))


class Email:
    def __init__(self, subject, template_name, **data):
        if config['Email']['Enabled']:
            # Initializes the email
            self.msg = MIMEMultipart('alternative')
            self.msg['From'] = config['Email']['Address']
            self.msg['Subject'] = subject

            # Parses the template
            template = templates.get_template(template_name + '.html')
            text = template.render(banner=config['Environment']['Address'] + '/static/banner.png', **data)
            self.msg.attach(MIMEText(text, 'html'))

    def send(self, *recipients):
        if config['Email']['Enabled']:
            # Connects to the server
            conn = SMTP(config['Email']['Server'], config['Email']['Port'])
            conn.ehlo()
            conn.starttls()
            conn.login(config['Email']['Login'], config['Email']['Password'])

            # Sends the email to everyone
            for recipient in recipients:
                del self.msg['To']
                self.msg['To'] = recipient
                conn.sendmail(config['Email']['Address'], recipient, self.msg.as_string())

            conn.close()

    def broadcast(self):
        # Sends the email to every user
        addresses = get_users_emails()
        self.send(*addresses)
        return len(addresses)
