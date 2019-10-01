package function

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"
)

func TestFailedLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago",
		httpmock.NewStringResponder(200, "wtv"),
	)

	httpmock.RegisterResponder(
		"POST",
		"https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc",
		httpmock.NewStringResponder(200, "wtv"),
	)

	_, err := Login("a", "b")

	assert.Equal(t, httpmock.GetTotalCallCount(), 2)
	info := httpmock.GetCallCountInfo()
	assert.Equal(
		t,
		info["POST https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago"],
		1,
	)
	assert.Equal(
		t,
		info["POST https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc"],
		1,
	)
	assert.Error(t, err, "invalid login")
}

func TestGoodLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	x1 := httpmock.NewStringResponse(302, `wtv`)
	x1.Header.Add("Location", "https://portalprepagos.cgd.pt/portalprepagos/private/home.seam")

	x2 := httpmock.NewStringResponse(200, `wtv`)
	x2.Header.Add("Set-Cookie", "JSESSIONID=123")
	x2.Header.Add("Set-Cookie", "SMSESSION=456")
	x2.Header.Add("Set-Cookie", "IGNORED=999")

	httpmock.RegisterResponder(
		"POST",
		"https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago",
		httpmock.NewStringResponder(200, "wtv"),
	)

	httpmock.RegisterResponder(
		"POST",
		"https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc",
		httpmock.ResponderFromResponse(x1),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://portalprepagos.cgd.pt/portalprepagos/private/home.seam",
		httpmock.ResponderFromResponse(x2),
	)

	token, err := Login("a", "b")
	assert.Equal(t, httpmock.GetTotalCallCount(), 3)
	info := httpmock.GetCallCountInfo()
	assert.Equal(
		t,
		info["POST https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago"],
		1,
	)
	assert.Equal(
		t,
		info["POST https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc"],
		1,
	)
	assert.Equal(
		t,
		info["GET https://portalprepagos.cgd.pt/portalprepagos/private/home.seam"],
		1,
	)
	assert.NilError(t, err)
	assert.Equal(t, token, "123#456")
}

func TestGetDataBadToken(t *testing.T) {
	data, err := GetData("")
	assert.Assert(t, data == nil)
	assert.Error(t, err, "invalid token")
}

func TestGetDataTokenExpired(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	x := httpmock.NewStringResponse(302, `wtv`)
	x.Header.Add("Location", "https://portalprepagos.cgd.pt/portalprepagos/private/home.seam")

	httpmock.RegisterResponder(
		"GET",
		"https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam",
		httpmock.ResponderFromResponse(x),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://portalprepagos.cgd.pt/portalprepagos/private/home.seam",
		httpmock.NewStringResponder(200, "wtv"),
	)

	data, err := GetData("a#b")
	assert.Assert(t, data == nil)
	assert.Error(t, err, "not logged in")
}

func TestGetData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	x := httpmock.NewStringResponse(
		200,
		string(helperLoadBytes(t, "getdata.html")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam",
		httpmock.ResponderFromResponse(x),
	)

	data, err := GetData("a#b")
	assert.NilError(t, err)
	assert.DeepEqual(t, data, &DataInfo{
		Balance: 121.12,
		PreviousExtracts: []string{
			"07/19",
			"06/19",
			"05/19",
			"04/19",
			"03/19",
			"02/19",
			"01/19",
			"12/18",
			"11/18",
			"10/18",
			"09/18",
		},
		CardNumber: "10122234853",
		ViewState:  "x:y",
		History: []Transaction{
			{
				Date:        "02-08-2019",
				ValueDate:   "01-08-2019",
				Description: "CASA DO PAO    4300-025 PORTO",
				DebitAmount: 5,
			},
			{
				Date:        "05-08-2019",
				ValueDate:   "02-08-2019",
				Description: "CASA DO BIFE   4000-334 PORTO",
				DebitAmount: 12.6,
			},
			{
				Date:        "05-08-2019",
				ValueDate:   "04-08-2019",
				Description: "PASTELARIA BOLA BERLIM SAO PEDRO DO SUL",
				DebitAmount: 9.45,
			},
			{
				Date:        "13-08-2019",
				ValueDate:   "12-08-2019",
				Description: "CASA DO BIFE      4000-334 PORTO",
				DebitAmount: 7.8,
			},
			{
				Date:        "16-08-2019",
				ValueDate:   "14-08-2019",
				Description: "FRANCESINHA  PORTO",
				DebitAmount: 11.3,
			},
			{
				Date:         "21-07-2019",
				ValueDate:    "21-07-2019",
				Description:  "CARREGAMENTO AUTOMATICO PRE-PAGO",
				CreditAmount: 93.23,
			},
		},
	})
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
