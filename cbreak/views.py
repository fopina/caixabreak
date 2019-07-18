from django.shortcuts import render
from django.http import JsonResponse

from Crypto.Cipher import AES
import os
import base64

from cbreak import Break, BreakException


def index(request):
    return render(
        request,
        'cbreak/index.html'
    )


def info(request):
    username = request.POST.get('username', None)
    password = request.POST.get('password', None)

    token = ''
    if not password:
        token = request.POST.get('token', '')
        if token:
            token = base64.b64decode(token)
            if len(token) == 32:
                encpwd = request.session.get('encpwd')
                if encpwd:
                    encpwd = base64.b64decode(encpwd)
                    secret = AES.new(token)
                    password = secret.decrypt(encpwd).rstrip('\0')

    response_data = {
        'success': False
    }

    if username and password:
        account = Break(username, password)
        try:
            account.login()
            a, b = account.balance_history()
            response_data['balance'] = a
            response_data['history'] = b
            response_data['success'] = True

            if not token:
                token = os.urandom(32)
                response_data['token'] = base64.b64encode(token)
                secret = AES.new(token)
                padpwd = password + ((AES.block_size - len(password)) % AES.block_size) * '\0'
                request.session['encpwd'] = base64.b64encode(secret.encrypt(padpwd))

        except BreakException as e:
            response_data['error'] = e.message
    else:
        response_data['error'] = 'Login not provided/expired'

    return JsonResponse(response_data)
