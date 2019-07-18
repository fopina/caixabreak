#!/bin/sh

if [ -n "${PG_PASSWORD_FILE}" ]; then
	export PG_PASSWORD=$(cat ${PG_PASSWORD_FILE})
fi

python /app/manage.py collectstatic --noinput
python /app/manage.py migrate
s6-svscan /s6
