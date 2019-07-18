import os

DEBUG = False
TEMPLATE_DEBUG = DEBUG

SECRET_KEY = os.getenv('DJANGO_SECRET_KEY')

if os.getenv('SQLITE_FILE'):
    DATABASES = {
        'default': {
            'ENGINE': 'django.db.backends.sqlite3',
            'NAME': os.getenv('SQLITE_FILE'),
        }
    }
else:
    DATABASES = {
        'default': {
            'ENGINE': 'django.db.backends.postgresql_psycopg2',
            'NAME': os.getenv('PG_DATABASE', 'break'),
            'USER': os.getenv('PG_USER', 'postgres'),
            'PASSWORD': os.getenv('PG_PASSWORD', ''),
            'HOST': os.getenv('PG_HOST', 'localhost'),
            'PORT': os.getenv('PG_PORT', ''),
        }
    }

ALLOWED_HOSTS = [ '*' ]
