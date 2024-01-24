FROM python:3.11-alpine

WORKDIR /cucinassistant

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY application/ application/
COPY config_sample.cfg config.cfg
COPY send_email.py send_email.py

CMD gunicorn -w 3 --threads 5 -b :80 application:app
