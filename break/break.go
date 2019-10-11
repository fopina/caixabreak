package function

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// https://www.cgd.pt/Particulares/Cartoes/Cartoes-Pre-pagos/Pages/Portal-pre-pagos.aspx
const urlData = "https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam"
const urlLogin = "https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago"
const urlLogin2 = "https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc"
const urlLoggedIn = "https://portalprepagos.cgd.pt/portalprepagos/private/home.seam"
const urlLogout = "https://portalprepagos.cgd.pt/portalprepagos/private/logout.seam?actionMethod=private%2Fhome.xhtml%3Aidentity.logout"
const loginTarget = "/portalprepagos/private/home.seam"
const userAgent = "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.1) Gecko/2008071615 Fedora/3.0.1-1.fc9 Firefox/3.0.1 Break/1.0"

const tokenCookie1 = "JSESSIONID"
const tokenCookie2 = "SMSESSION"

var /* const */ dayRows = regexp.MustCompile(`(?s)<tr class="dia-(.+?)\b.*?">(.*?)</tr>`)
var /* const */ dayCols = regexp.MustCompile(`(?s)<td .*?class="(.*?)">(.*?)</td>`)
var /* const */ balanceRE = regexp.MustCompile(`(?s)<p class="ref">Saldo disponível\s*</p>\s+<p class="valor">\s*<label>\s*(\d+,\d\d)</label><label class="marginBetween">\s*EUR</label>`)
var /* const */ entriesRE = regexp.MustCompile(
	`(?s)<tr onmouseover="JavaScript:this.className=\'row_Over\'" onmouseout="JavaScript:this.className=\'.+?\';" class=".+?">` +
		`\s*<td width="15%" class="texttable col_fst">(.+?)\s*</td>\s+<td width="15%" class="texttable col_med">(.+?)\s*</td>` +
		`\s+<td class="texttable col_med">(.+?)\s*</td>\s*<td class="texttable col_med alignrighttext">(.+?)\s*</td>\s*` +
		`<td class="texttable  alignrighttext col_last">(.+?)\s*</td>\s*</tr>`,
)
var /* const */ monthOptionsRE = regexp.MustCompile(`<option value="(\d\d/\d\d)".*?>\d\d-\d\d\d\d</option>`)
var /* const */ viewStateRE = regexp.MustCompile(
	`<input type="hidden" name="javax.faces.ViewState" id="javax.faces.ViewState" value="(.*?)" autocomplete="off" />`,
)
var /* const */ cardNumberRE = regexp.MustCompile(
	`<select id="consultaMovimentosCartoesPrePagos:selectedCard" name="consultaMovimentosCartoesPrePagos:selectedCard" class="componentsComboBox" size="1" onchange="app.post\(this\);window.focus\(\);">\s*<option value="(\d+)" selected="selected">`,
)

// Transaction holds transaction info
type Transaction struct {
	Date,
	ValueDate,
	Description string
	DebitAmount,
	CreditAmount float64
}

// DataInfo holds information per page - balance and history
type DataInfo struct {
	Balance          float64
	PreviousExtracts []string
	CardNumber,
	ViewState string
	History []Transaction
}

// UnauthorizedError error thrown when bad credentials are provided (user and password or invalid/expired token)
type UnauthorizedError struct {
	s string
}

func (e *UnauthorizedError) Error() string {
	return e.s
}

func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

// GetToken extract token from cookiejar
func GetToken(client *http.Client) string {
	u, err := url.Parse(urlLoggedIn)
	if err != nil {
		log.Fatal(err)
	}

	token := make([]string, 2)

	for _, c := range client.Jar.Cookies(u) {
		if c.Name == tokenCookie1 {
			token[0] = c.Value
		}
		if c.Name == tokenCookie2 {
			token[1] = c.Value
		}
	}

	return strings.Join(token, "#")
}

// Login log in to Sigarra
func Login(username, password string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Jar: jar}

	form := url.Values{
		"USERNAME":    {username},
		"login_btn_1": {"OK"},
	}

	req, err := newRequest("POST", urlLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.Request.URL.String() != urlLogin {
		return "", fmt.Errorf("unexpect redirect, maintenance?")
	}

	form = url.Values{
		"target":           {loginTarget},
		"username":         {"PPP" + username},
		"userInput":        {username},
		"passwordInput":    {"*****"},
		"loginForm:submit": {"Entrar"},
		"password":         {password},
	}
	req, err = newRequest("POST", urlLogin2, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.Request.URL.String() != urlLoggedIn {
		return "", &UnauthorizedError{"invalid login"}
	}

	return GetToken(client), nil
}

// Logout terminate session in Sigarra
func Logout(tokenString string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client := &http.Client{Jar: jar}

	err = setupJar(tokenString, jar)
	if err != nil {
		return err
	}

	form := url.Values{
		"p_address": {"loginAddress"},
	}
	req, err := newRequest("POST", urlLogout, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extractData(client *http.Client, req *http.Request) (*DataInfo, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Request.URL.String() != urlData {
		return nil, &UnauthorizedError{"not logged in"}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)

	data := &DataInfo{}

	for _, match := range balanceRE.FindAllStringSubmatch(html, -1) {
		data.Balance, _ = parseThisFloat(match[1])
	}

	for _, match := range viewStateRE.FindAllStringSubmatch(html, -1) {
		data.ViewState = match[1]
	}

	for _, match := range cardNumberRE.FindAllStringSubmatch(html, -1) {
		data.CardNumber = match[1]
	}

	entries := entriesRE.FindAllStringSubmatch(html, -1)
	data.History = make([]Transaction, len(entries))
	for i, match := range entries {
		data.History[i].Date = match[1]
		data.History[i].ValueDate = match[2]
		data.History[i].Description = match[3]
		data.History[i].DebitAmount, err = parseThisFloat(match[4])
		if err != nil {
			data.History[i].DebitAmount = 0
		}
		data.History[i].CreditAmount, err = parseThisFloat(match[5])
		if err != nil {
			data.History[i].CreditAmount = 0
		}
	}

	entries = monthOptionsRE.FindAllStringSubmatch(html, -1)
	data.PreviousExtracts = make([]string, len(entries))
	for i, match := range entries {
		data.PreviousExtracts[i] = match[1]
	}

	return data, nil
}

// GetData retrieve the attendance data from Sigarra
func GetData(tokenString string) (*DataInfo, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}

	err = setupJar(tokenString, jar)
	if err != nil {
		return nil, err
	}

	req, err := newRequest("GET", urlData, nil)
	if err != nil {
		return nil, err
	}

	return extractData(client, req)
}

// GetDataForMonth retrieve the attendance data from Sigarra for specific MM/YY
func GetDataForMonth(tokenString, viewState, cardNumber, monthYear string) (*DataInfo, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}

	err = setupJar(tokenString, jar)
	if err != nil {
		return nil, err
	}

	// not including viewState seems to be a valid viewState, nice!
	form := url.Values{
		"consultaMovimentosCartoesPrePagos":                  {"consultaMovimentosCartoesPrePagos"},
		"consultaMovimentosCartoesPrePagos:ignoreFieldsComp": {""},
		"consultaMovimentosCartoesPrePagos:selectedCard":     {cardNumber},
		"consultaMovimentosCartoesPrePagos:extractDates":     {monthYear},
		"javax.faces.ViewState":                              {viewState},
	}
	req, err := newRequest("POST", urlData, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return extractData(client, req)
}

func setupJar(tokenString string, jar *cookiejar.Jar) error {
	u, err := url.Parse(urlData)
	if err != nil {
		return err
	}
	token := strings.Split(tokenString, "#")
	if len(token) != 2 {
		return &UnauthorizedError{"invalid token"}
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   tokenCookie1,
		Value:  token[0],
		Path:   "/",
		Domain: u.Hostname(),
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:   tokenCookie2,
		Value:  token[1],
		Path:   "/",
		Domain: u.Hostname(),
	}
	cookies = append(cookies, cookie)
	jar.SetCookies(u, cookies)
	return nil
}

func parseThisFloat(floatString string) (float64, error) {
	return strconv.ParseFloat(
		strings.ReplaceAll(
			strings.ReplaceAll(floatString, ".", ""),
			",", ".",
		),
		64,
	)
}
