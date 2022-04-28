package tokens_test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/proactiongo/pagocore"
	"github.com/proactiongo/pagocore/i18n"
	"github.com/proactiongo/pagocore/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseAccessToken(t *testing.T) {
	initial := pagocore.Opt.JWTPassword
	pagocore.Opt.JWTPassword = []byte("test")

	valid := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIyOTI3ZTgxMC1lOTQ4LTRiN2UtYjBkMi1lZThhMWYxNzFkMGYiLCJleHAiOjU2MjUxNTI0OTgsImp0aSI6IjdjZDY3NzQyLTNlMjMtNDdmZS04OTQzLWZhZWI1NjM4NzNkMSIsInJsIjo1MDAsIm5tIjoiRGV2IFRlc3QgVGVhY2hlciIsInJ0aSI6ImNmZGRlNzVkLTQzY2ItNGI5Ni1iZjdjLTk4ODBiMzY5Zjk1OSIsImxuZyI6InJ1In0.c-gTREm9retM50F3ZZqStK14qKINrOsmoV_VnNScjps",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIyOTI3ZTgxMC1lOTQ4LTRiN2UtYjBkMi1lZThhMWYxNzFkMGYiLCJleHAiOjU2MjUxNTI0OTgsImp0aSI6IjdjZDY3NzQyLTNlMjMtNDdmZS04OTQzLWZhZWI1NjM4NzNkMSIsInJsIjoxLCJubSI6IkRldiBUZXN0IFVzZXIiLCJydGkiOiJjZmRkZTc1ZC00M2NiLTRiOTYtYmY3Yy05ODgwYjM2OWY5NTkiLCJsbmciOiJydSJ9.Dt8FJme4S1zM8A35HkC5Xm_IbjFPtnguRrYQ_kD9BOc",
	}

	invalid := []string{
		// invalid password
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Y2Xw2UjBZHPfGOfbm3ryZlIpdfcESt1ewcLCaDbSgv0",
		// no IDs
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjU2MjUxNTI0OTgsInJsIjoxLCJubSI6IkRldiBUZXN0IFVzZXIiLCJydGkiOiJjZmRkZTc1ZC00M2NiLTRiOTYtYmY3Yy05ODgwYjM2OWY5NTkiLCJsbmciOiJydSJ9.WFv4QeqQgClencOgmH5ScACQECw6SdZvS_-2O-vz1UQ",
		// expired
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIyOTI3ZTgxMC1lOTQ4LTRiN2UtYjBkMi1lZThhMWYxNzFkMGYiLCJleHAiOjEwMCwianRpIjoiN2NkNjc3NDItM2UyMy00N2ZlLTg5NDMtZmFlYjU2Mzg3M2QxIiwicmwiOjEsIm5tIjoiRGV2IFRlc3QgVXNlciIsInJ0aSI6ImNmZGRlNzVkLTQzY2ItNGI5Ni1iZjdjLTk4ODBiMzY5Zjk1OSIsImxuZyI6InJ1In0.rZ6Q2sR8qRs6aXysJKTFsxYtxHcSogwKciGBUcoFuH0",
		// invalid method
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IldURiJ9.eyJ1aWQiOiIyOTI3ZTgxMC1lOTQ4LTRiN2UtYjBkMi1lZThhMWYxNzFkMGYiLCJleHAiOjEwMCwianRpIjoiN2NkNjc3NDItM2UyMy00N2ZlLTg5NDMtZmFlYjU2Mzg3M2QxIiwicmwiOjEsIm5tIjoiRGV2IFRlc3QgVXNlciIsInJ0aSI6ImNmZGRlNzVkLTQzY2ItNGI5Ni1iZjdjLTk4ODBiMzY5Zjk1OSIsImxuZyI6InJ1In0.Jxqfk-g8ZCXp_b4p8BuSi6QId0V625gLsUnwcsPeGdQ",
		// invalid rti
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IldUIn0.eyJ1aWQiOiIyOTI3ZTgxMC1lOTQ4LTRiN2UtYjBkMi1lZThhMWYxNzFkMGYiLCJleHAiOjEwMCwianRpIjoiN2NkNjc3NDItM2UyMy00N2ZlLTg5NDMtZmFlYjU2Mzg3M2QxIiwicmwiOjEsIm5tIjoiRGV2IFRlc3QgVXNlciIsInJ0aSI6ImludmFsaWQiLCJsbmciOiJydSJ9.y_pm-1w4aLn_9VxT6F2caasOxtMVkUxn1e2ARotiulY",
	}

	for _, v := range valid {
		_, err := tokens.ParseAccessToken(v)
		assert.NoError(t, err)
	}

	for _, v := range invalid {
		_, err := tokens.ParseAccessToken(v)
		assert.ErrorIs(t, err, pagocore.ErrTokenInvalid)
	}

	pagocore.Opt.JWTPassword = initial
}

func TestAccessTokenClaims_Valid(t *testing.T) {
	initial := pagocore.Opt.ServiceName
	pagocore.Opt.ServiceName = "test"

	valid := &tokens.AccessTokenClaims{
		RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		Role:           10,
		Name:           "User Name",
		Language:       i18n.LangRu,
	}
	valid.TokenClaimsDft = tokens.TokenClaimsDft{
		UserID:          "03a4e59c-fb22-4bfa-8739-8062bcdd2005",
		ServicesAllowed: []string{"test1", "test"},
		StandardClaims: jwt.StandardClaims{
			Id:        "0bf97df4-6246-4809-bdf7-e8d993668283",
			ExpiresAt: time.Now().Unix() + int64(10*time.Minute.Seconds()),
		},
	}
	assert.NoError(t, valid.Valid())

	pagocore.Opt.ServiceName = ""
	assert.Error(t, valid.Valid())

	pagocore.Opt.ServiceName = "test"

	invalid := []*tokens.AccessTokenClaims{
		{
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
			Role:           10,
			Name:           "User Name",
			Language:       i18n.LangRu,
			TokenClaimsDft: tokens.TokenClaimsDft{
				UserID:          "03a4e59c-fb22-4bfa-8739-8062bcdd2005",
				ServicesAllowed: []string{"test2", "test3"},
				StandardClaims: jwt.StandardClaims{
					Id:        "0bf97df4-6246-4809-bdf7-e8d993668283",
					ExpiresAt: time.Now().Unix() + int64(10*time.Minute.Seconds()),
				},
			},
		},
		{},
		{Role: 10, RefreshTokenID: ""},
		{
			TokenClaimsDft: tokens.TokenClaimsDft{},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: tokens.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: tokens.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
				StandardClaims: jwt.StandardClaims{
					Id:        "dd5d2ced-e168-4794-822e-f13fe952dddd",
					ExpiresAt: 100,
				},
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: tokens.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
				StandardClaims: jwt.StandardClaims{
					Id:        "dd5d2ced-e168-4794-822e-f13fe952dddd",
					ExpiresAt: time.Now().Unix() + int64(10*time.Minute.Seconds()),
					NotBefore: time.Now().Unix() + int64(5*time.Minute.Seconds()),
				},
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
	}

	for _, claims := range invalid {
		assert.Error(t, claims.Valid())
	}

	pagocore.Opt.ServiceName = initial
}
