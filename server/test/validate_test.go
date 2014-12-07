package api_test

import (
	"../types"
	. "gopkg.in/check.v1"
)

// Types specific for testing validation
type handleProposal struct {
	Handle string `validate:"handle"`
}

type emailProposal struct {
	Email string `validate:"email"`
}

type passwordProposal struct {
	Password string `validate:"password"`
}

func (s *TestSuite) TestValidateBadHandle(c *C) {
	v := types.NewValidator()
	badHandles := []handleProposal{
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
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, NotNil, Commentf("Index %d has no error", i))
	}
}

func (s *TestSuite) TestValidateGoodHandle(c *C) {
	v := types.NewValidator()
	goodHandles := []handleProposal{
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
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, IsNil, Commentf("Index %d: %s", i, result))
	}
}

func (s *TestSuite) TestValidateBadEmail(c *C) {
	v := types.NewValidator()
	badEmails := []emailProposal{
		{Email: "joe@usa"},
		{Email: "sonin@.com"},
		{Email: "samurai@@gmail.com"},
		{Email: "mocha@latte@io.com"},
		{Email: "@microsoft@tech.it"},
		{Email: "treehouse"},
		{Email: "football@nfl..com"},
		{Email: "toolbar@api.microsoft..com"},
		{Email: "email@password.io@email.com"},
		{Email: "erin@apple,com"},
	}

	for i, proposal := range badEmails {
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, NotNil, Commentf("Index %d has no error", i))
	}
}

func (s *TestSuite) TestValidateGoodEmail(c *C) {
	v := types.NewValidator()
	goodEmails := []emailProposal{
		{Email: "edward@cherami.io"},
		{Email: "tester@microsoft.com"},
		{Email: "t234342@gmail.com"},
		{Email: "uncle-bob@biz.info"},
		{Email: "mkdir@derp.museum"},
		{Email: "go-lang@static.tk"},
		{Email: "gopher@golang.google.co"},
		{Email: "zqk@test.merra.is"},
		{Email: "bigballer@swag.gov"},
		{Email: "tek421@sub.domain.biz"},
		{Email: "4213@4uandme.com"},
	}
	for i, proposal := range goodEmails {
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, IsNil, Commentf("Index %d: %s", i, result))
	}
}

func (s *TestSuite) TestValidateBadPassword(c *C) {
	v := types.NewValidator()
	badPasswords := []passwordProposal{
		{Password: ""},
		{Password: "aaaaaaa"},
		{Password: "VLeHkciByWBXNnaExhMA6QKwioybgEZCkEj9YzyhwvbofKTejj1"},
	}
	for i, proposal := range badPasswords {
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, NotNil, Commentf("Index %d has no error", i))
	}
}

func (s *TestSuite) TestValidateGoodPassword(c *C) {
	v := types.NewValidator()
	goodPasswords := []passwordProposal{
		{Password: "6FxDW9ws"},
		{Password: "(y/&63N79;,6{36bp^7x=(7>8CZi"},
		{Password: "iMpbnuVZadZKCYbTwoDbgmLfTNUvQDzRpdBxfWrbZCUHXkzEBx"},
	}
	for i, proposal := range goodPasswords {
		result := v.ValidateAndTag(proposal, "json")
		c.Check(result, IsNil, Commentf("Index %d: %s", i, result))
	}
}
