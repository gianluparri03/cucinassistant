FROM python:3.8-alpine

WORKDIR /cucinassistant
COPY application/ application/

COPY requirements.txt .
RUN apk add gcc musl-dev mariadb-connector-c-dev curl
RUN pip install -r requirements.txt

COPY run.py .
ENV PRODUCTION=true
CMD gunicorn -w 1 --threads 5 -b :80 application:app
