FROM python:3.11-alpine

WORKDIR /cucinassistant

RUN apk add gcc musl-dev mariadb-connector-c-dev curl

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY cucinassistant/ cucinassistant/
COPY send_email.py send_email.py

CMD gunicorn -w 3 --threads 5 -b :80 cucinassistant:app
