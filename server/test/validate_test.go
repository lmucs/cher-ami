package api_test

import (
	"../types"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestValidateGoodHandle(c *C) {
	v := types.NewValidator()
	goodHandles := []types.SignupProposal{
		{Handle: "John400"},
		{Handle: "タロウ"},
		{Handle: "やまだ"},
		{Handle: "山田"},
		{Handle: "先生"},
		{Handle: "ыхаыл"},
		{Handle: "Θεοκλεια"},
		{Handle: "आकाङ्क्षा"},
		{Handle: "אַבְרָהָם"},
		{Handle: "മലയാളം"},
		{Handle: "상"},
	}

	for i, proposal := range goodHandles {
		result := v.Validate(proposal)
		c.Check(result, IsNil, Commentf("Index %d: %s", i, result))
	}
}

func (s *TestSuite) TestValidateBadHandle(c *C) {
	v := types.NewValidator()
	badHandles := []types.SignupProposal{
		{Handle: ""},
		{Handle: "400John"},
		{Handle: "タロウタウタウタウウタウタウタウタウタロウ"},
		{Handle: "山田κλειαയാള상आकाङ्GGQQ"},
		{Handle: "山田:山田"},
		{Handle: "#ыхаыл"},
		{Handle: "@#$%^&*(****%^&*("},
		{Handle: "&"},
		{Handle: "2g"},
		{Handle: "や##ま$だ"},
	}

	for i, proposal := range badHandles {
		result := v.Validate(proposal)
		c.Check(result, NotNil, Commentf("Index %d has no error", i))
	}
}
