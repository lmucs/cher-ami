package api_test

import (
	"../types"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestValidateHandle(c *C) {
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
