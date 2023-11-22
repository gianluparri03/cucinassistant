FROM python:3.11-alpine

WORKDIR /cucinassistant
COPY application/ application/

COPY requirements.txt .
RUN apk add curl
RUN pip install -r requirements.txt

COPY run.py .
ENV PRODUCTION=true
CMD gunicorn -w 3 --threads 5 -b :80 application:app
