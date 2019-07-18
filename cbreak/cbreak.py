# -*- coding: utf-8 -*-
import re, requests

requests.packages.urllib3.disable_warnings()


class BreakException(Exception):
    pass


class Break(object):
    def __init__(self, user, password):
        self._session = requests.Session()
        self._user = user
        self._password = password

    def login(self):
        # url = 'https://www.cgd.pt/particulares/gerir-dia-a-dia/cartoes/Portal-pre-pagos/Pages/Portal-pre-pagos.aspx'
        payload = {
            'USERNAME': self._user,
            'login_btn_1': 'OK',
        }
        r = self._session.post(
            'https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago',
            data=payload
            )

        if r.url != 'https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago':
            # redirected somewhere else?? Can it happen?
            # better safe than sorry
            raise BreakException('Failed to login')

        payload = {
            'target': '/portalprepagos/private/home.seam',
            'username': 'PPP%s' % self._user,
            'userInput': self._user,
            'passwordInput': '*****',
            'loginForm:submit': 'Entrar',
            'password': self._password,
        }
        r = self._session.post(
            'https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc',
            data = payload
            )

        # login error
        # ToDo - parse error message
        if r.url != 'https://portalprepagos.cgd.pt/portalprepagos/private/home.seam':
            raise BreakException('Failed to login')

    def balance_history(self):
        r = self._session.get('https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam')
        m = re.search(
            u'<p class="ref">Saldo disponível\s*</p>\s+<p class="valor">'
            u'<label>\s*(\d+,\d\d)</label><label class="marginBetween">\s*EUR</label>',
            r.text,
            re.S
        )
        bal = float(m.group(1).replace(',', '.'))
        m = re.findall(
            u'<tr onmouseover="JavaScript:this.className=\'row_Over\'" onmouseout="JavaScript:this.className=\'.+?\';" class=".+?">'
            u'\s*<td width="15%" class="texttable col_fst">(.+?)\s*</td>\s+<td width="15%" class="texttable col_med">(.+?)\s*</td>'
            u'\s+<td class="texttable col_med">(.+?)\s*</td>\s*<td class="texttable col_med alignrighttext">(.+?)\s*</td>\s*'
            u'<td class="texttable  alignrighttext col_last">(.+?)\s*</td>\s*</tr>',
            r.text,
            re.S)

        mov = []

        for n in m:
            nn = list(n[:3])
            if n[3] == '&nbsp;':
                nn.append(float(n[4].replace(',','.')))
            else:
                nn.append(-float(n[3].replace(',','.')))
            mov.append(nn)

        return (bal, mov)