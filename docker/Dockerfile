FROM python:2.7.13-alpine3.6 as builder

RUN apk add --no-cache postgresql-dev g++ make

ADD requirements*.txt /
RUN pip wheel -w /wheels -r requirements_extra.txt 

FROM python:2.7.13-alpine3.6

RUN apk add --no-cache postgresql-dev nginx s6

ADD . /app/

COPY --from=builder /wheels /wheels
RUN pip install --no-index --find-links=/wheels -r /app/requirements_extra.txt

RUN mv /app/docker/local_settings.py /app/caixabreak/local_settings.py
RUN mv /app/docker/nginx.conf /etc/nginx/nginx.conf
RUN mv /app/docker/s6 /s6

RUN mkdir -p /run/nginx/

EXPOSE 8080
WORKDIR /app
CMD "/app/docker/start.sh"
