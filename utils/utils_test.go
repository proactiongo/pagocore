package utils_test

import (
	"github.com/proactiongo/pagocore/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

func TestExtractBearerToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1", nil)
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set("Authorization", "Bearer TEST_TOKEN")

	token, err := utils.ExtractBearerToken(req)
	assert.NoError(t, err)
	assert.Equal(t, "TEST_TOKEN", token)

	req, err = http.NewRequest(http.MethodPost, "http://127.0.0.1", nil)
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set("Authorization", "TEST_TOKEN")

	_, err = utils.ExtractBearerToken(req)
	assert.Error(t, err)

	req, err = http.NewRequest(http.MethodPost, "http://127.0.0.1", nil)
	if !assert.NoError(t, err) {
		return
	}

	_, err = utils.ExtractBearerToken(req)
	assert.Error(t, err)
}

func TestValidateUUID(t *testing.T) {
	valid := []string{
		"8c03c25b-e1cc-4565-a41d-caf8900b0f7f",
		"D2B2660E-47ED-4309-9DC9-2DDBEE497D14",
		utils.GenerateUUID(),
		utils.GenerateUUID(),
	}
	invalid := []string{
		"6a41b89a-zZz-4c61-a19a-304533ee6a6c",
		"6a41b89a-4c61-a19a-304533ee6a6c",
		"some other string",
		"",
	}

	for _, uuid := range valid {
		if !utils.ValidateUUID(uuid) {
			t.Error("valid UUID failed", uuid)
		}
	}

	for _, uuid := range invalid {
		if utils.ValidateUUID(uuid) {
			t.Error("invalid UUID passed", uuid)
		}
	}
}

func TestFilterEmail(t *testing.T) {
	emails := map[string]string{
		" \njohn_doe@example.com ":          "john_doe@example.com",
		" John Doe <john_doe@example.com> ": "john_doe@example.com",

		"   ": "",
	}
	for input, expected := range emails {
		res := utils.FilterEmail(input)
		assert.Equal(t, expected, res)
	}
}

func TestValidateEmail(t *testing.T) {
	valid := []string{
		" Some Name <some-email.1@example.com> ",
		"ктовообщепользуетсярусскими@адресами.рф",
	}
	invalid := []string{
		"this is not an email @",
		"123",
		"",
	}

	for _, val := range valid {
		ok := utils.ValidateEmail(val)
		assert.Equal(t, true, ok)
	}

	for _, val := range invalid {
		ok := utils.ValidateEmail(val)
		assert.Equal(t, false, ok)
	}
}

func TestValidatePhone(t *testing.T) {
	valid := []string{
		"+7123-45678 90",
		"+(123) 456-12-45",
		"89213456789",
	}
	invalid := []string{
		"this is not a phone number",
		"123",
		"+",
		"",
	}

	for _, val := range valid {
		if !utils.ValidatePhone(val) {
			t.Error("unexpected phone validation failure (", val, ")")
		}
	}

	for _, val := range invalid {
		if utils.ValidatePhone(val) {
			t.Error("unexpected phone validation pass (", val, ")")
		}
	}
}

func TestJSONConvert(t *testing.T) {
	sourceValid := map[string]interface{}{
		"teacher_id": "5c1056c5-292d-43dd-bf3c-d0bd0e2b7064",
		"student_id": "3edfe19d-143c-4cd4-81b9-c368ebc6da7e",
	}
	targetValid := map[string]string{}

	err := utils.JSONConvert(sourceValid, &targetValid)
	assert.NoError(t, err)

	assert.Equal(t, sourceValid["teacher_id"].(string), targetValid["teacher_id"])
	assert.Equal(t, sourceValid["student_id"].(string), targetValid["student_id"])

	err = utils.JSONConvert("invalid value", &targetValid)
	assert.Error(t, err)
}

func TestValidateURL(t *testing.T) {
	{
		t.Log("Testing valid URLs")
		valid := []string{
			"https://example.com",
			"http://example.com?test=1",
			"//example.com",
		}

		for i, s := range valid {
			assert.Equal(t, true, utils.ValidateURL(s), i)
		}
	}

	{
		t.Log("Testing invalid URLs")
		invalid := []string{
			"what is it?",
			"42 is not an answer",
			"",
			"example.com/test",
		}

		for i, s := range invalid {
			assert.Equal(t, false, utils.ValidateURL(s), i)
		}
	}
}

func TestFilterStrings(t *testing.T) {
	input := []string{
		"",
		"   ",
		"\n\n",
		" TEST1 ",
		"  TEST2  ",
	}
	expected := []string{
		"TEST1",
		"TEST2",
	}
	output := utils.FilterStrings(input)

	assert.EqualValues(t, expected, output)

	nilOut := utils.FilterStrings(nil)
	assert.NotNil(t, nilOut)
	assert.Equal(t, 0, len(nilOut))
}

func TestGenerateRandomString(t *testing.T) {
	symbols := "abcdef0123"
	count := 6

	str := utils.GenerateRandomString(symbols, uint(count))
	assert.Equal(t, count, len(str))

	rx, err := regexp.Compile("^[" + symbols + "]*$")
	if !assert.NoError(t, err) {
		return
	}

	match := rx.Match([]byte(str))
	assert.Equal(t, true, match)
}

func TestNormalizePhone(t *testing.T) {
	values := map[string]string{
		"+7 (900) 123-45-67": "+79001234567",
		"(900) 123-45-67":    "+79001234567",
		"8900 123-45-67":     "+79001234567",
		"89004567890":        "+79004567890",
		"9004567890":         "+79004567890",
	}

	for inp, expected := range values {
		assert.Equal(t, expected, utils.NormalizePhone(inp))
	}
}
