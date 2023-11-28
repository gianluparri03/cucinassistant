FROM python:3.11-alpine

WORKDIR /cucinassistant
COPY application/ application/

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY run.py .
RUN mkdir data/

ENV PRODUCTION=true
CMD export SECRET=`python -c 'import uuid; print(uuid.uuid4())'` && \
    gunicorn -w 3 --threads 5 -b :80 application:app
