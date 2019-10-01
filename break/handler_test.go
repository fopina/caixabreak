package function

import (
	"testing"

	"github.com/jarcoal/httpmock"
	handler "github.com/openfaas-incubator/go-function-sdk"
	"gotest.tools/assert"
)

func _TestInvalidLoginRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam",
		httpmock.NewStringResponder(200, "wtv"),
	)

	req := handler.Request{
		Body:        []byte("oi"),
		QueryString: "/",
	}

	res, err := Handle(req)

	assert.Equal(t, res.StatusCode, 401)
	assert.Equal(
		t, string(res.Body),
		`{"Token":"","Error":"invalid login","Days":null}`,
	)
	assert.Equal(t, httpmock.GetTotalCallCount(), 1)
	info := httpmock.GetCallCountInfo()
	assert.Equal(
		t,
		info["POST https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam"],
		1,
	)
	assert.NilError(t, err)
}

func _TestLoginRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://sigarra.up.pt/feup/pt/vld_validacao.validacao",
		httpmock.NewStringResponder(
			200,
			`<meta http-equiv="Refresh" content="0;url=https://sigarra.up.pt/feup/pt/ASSD_TLP_GERAL.FUNC_VIEW">`,
		),
	)
	x := httpmock.NewStringResponse(
		200,
		string(helperLoadBytes(t, "getdata.html")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://sigarra.up.pt/feup/pt/ASSD_TLP_GERAL.FUNC_VIEW",
		httpmock.ResponderFromResponse(x),
	)

	req := handler.Request{
		Body:        []byte("oi"),
		QueryString: "/",
	}

	res, err := Handle(req)

	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(
		t, string(res.Body),
		`{"Token":"#","Error":"","Days":[{"Type":"normal","Date":"2019-07-01","Balance":"0:05","BalanceAccrual":"0:05","Unjustified":"0:00","UnjustifiedAccrual":"0:00","MorningIn":"09:41","MorningOut":"13:02","AfternoonIn":"13:55","AfternoonOut":"16:46"},{"Type":"normal","Date":"2019-07-02","Balance":"0:16","BalanceAccrual":"0:21","Unjustified":"0:00","UnjustifiedAccrual":"0:00","MorningIn":"09:44","MorningOut":"12:40","AfternoonIn":"13:41","AfternoonOut":"17:01"},{"Type":"normal","Date":"2019-07-03","Balance":"0:20","BalanceAccrual":"0:41","Unjustified":"0:00","UnjustifiedAccrual":"0:00","MorningIn":"09:20","MorningOut":"13:02","AfternoonIn":"13:56","AfternoonOut":"16:40"},{"Type":"normal","Date":"2019-07-04","Balance":"0:18","BalanceAccrual":"0:59","Unjustified":"2:30","UnjustifiedAccrual":"2:30","MorningIn":"09:47","MorningOut":"12:32","AfternoonIn":"13:37","AfternoonOut":"17:10"},{"Type":"actual","Date":"2019-07-05","Balance":"-3:05","BalanceAccrual":"0:59","Unjustified":"0:00","UnjustifiedAccrual":"0:00","MorningIn":"09:47","MorningOut":"---","AfternoonIn":"---","AfternoonOut":"---"}]}`,
	)
	assert.Equal(t, httpmock.GetTotalCallCount(), 2)
	info := httpmock.GetCallCountInfo()
	assert.Equal(
		t,
		info["POST https://sigarra.up.pt/feup/pt/vld_validacao.validacao"],
		1,
	)
	assert.Equal(
		t,
		info["GET https://sigarra.up.pt/feup/pt/ASSD_TLP_GERAL.FUNC_VIEW"],
		1,
	)
	assert.NilError(t, err)
}
