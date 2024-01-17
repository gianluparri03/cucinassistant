FROM python:3.11-alpine

WORKDIR /cucinassistant

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY application/ application/
COPY run.py .

ENV PRODUCTION=true
RUN export SECRET=`python -c 'import secrets; print(secrets.token_hex())'`

CMD gunicorn -w 3 --threads 5 -b :80 application:app
