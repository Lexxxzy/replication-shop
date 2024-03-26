FROM locustio/locust:latest as benchmark

COPY locustfiles/requirements.txt /tmp/requirements.txt

RUN pip install -r /tmp/requirements.txt

ENTRYPOINT ["locust"]
