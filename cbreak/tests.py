from django.test import TestCase
from django.conf import settings
from django.test.client import Client
from django.core.urlresolvers import reverse
import json

from cbreak import Break, BreakException


class BreakTest(TestCase):
    def test_invalid_login(self):
        c = Break('0000000','12345')
        self.assertRaises(BreakException, c.login)

    def test_login(self):
        c = Break(settings.TEST_CBREAK_LOGIN, settings.TEST_CBREAK_PASSWORD)
        c.login()

    def test_history(self):
        c = Break(settings.TEST_CBREAK_LOGIN, settings.TEST_CBREAK_PASSWORD)
        c.login()

        balmov = c.balance_history()

        # two values returned
        self.assertEqual(len(balmov), 2)

        bal, mov = balmov

        # random assert to verify float conversion worked
        self.assertGreaterEqual(bal, 0)

        for m in mov:
            self.assertTrue(valid_date(m[0]))
            self.assertTrue(valid_date(m[1]))
            if m[2].startswith('CARREGAMENTO'):
                self.assertGreaterEqual(m[3], 0)
            else:
                self.assertLessEqual(m[3], 0)


class BreakWebTest(TestCase):
    def setUp(self):
        self.c = Client()

    def test_login_fail(self):
        res = self.c.post(reverse('cbreak:info'), {
            'username': 'invalid',
            'password': 'invalid',
            })
        r = json.loads(res.content)

        self.assertFalse(r['success'])
        self.assertEqual(r.get('error'), 'Failed to login')
        self.assertFalse('history' in r)
        self.assertFalse('balance' in r)

    def test_info(self):
        res = self.c.post(reverse('cbreak:info'), {
            'username': settings.TEST_CBREAK_LOGIN,
            'password': settings.TEST_CBREAK_PASSWORD,
            })
        r = json.loads(res.content)

        self.assertTrue(r['success'])
        self.assertFalse('error' in r)
        self.assertTrue('history' in r)
        self.assertTrue('balance' in r)

        # random assert to verify float conversion worked
        self.assertGreaterEqual(r['balance'], 0)

        for m in r['history']:
            self.assertTrue(valid_date(m[0]))
            self.assertTrue(valid_date(m[1]))
            if m[2].startswith('CARREGAMENTO'):
                self.assertGreaterEqual(m[3], 0)
            else:
                self.assertLessEqual(m[3], 0)

    def test_session(self):
        res = self.c.post(reverse('cbreak:info'), {
            'username': settings.TEST_CBREAK_LOGIN,
            'password': settings.TEST_CBREAK_PASSWORD,
            })
        r = json.loads(res.content)

        self.assertTrue(r['success'])
        self.assertFalse('error' in r)
        self.assertGreaterEqual(r['balance'], 0)
        self.assertTrue('token' in r)

        token = r['token']

        res = self.c.post(reverse('cbreak:info'), {
            'username': settings.TEST_CBREAK_LOGIN,
            'token': token
            })
        r = json.loads(res.content)

        self.assertTrue(r['success'])
        self.assertFalse('error' in r)
        self.assertGreaterEqual(r['balance'], 0)
        self.assertFalse('token' in r)


def valid_date(date_string):
    import re
    return re.match(r'\d\d-\d\d-\d\d\d\d', date_string) is not None
