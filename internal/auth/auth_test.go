package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashAndCheck(t *testing.T) {
	const password1 = "assword"
	const password2 = "a$$word"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTs(t *testing.T) {
	userID := uuid.New()
	root_secret := "MySecret"

	tests := []struct {
		name          string
		secret        string
		expiresIn     time.Duration
		wantErr       bool
		malformJWTStr bool
	}{
		{
			name:          "Valid",
			secret:        root_secret,
			expiresIn:     time.Second * 10,
			wantErr:       false,
			malformJWTStr: false,
		},
		{
			name:          "Expired",
			secret:        root_secret,
			expiresIn:     time.Minute * -10,
			wantErr:       true,
			malformJWTStr: false,
		},
		{
			name:          "Secret doesn't match",
			secret:        "different secret",
			expiresIn:     time.Second * 10,
			wantErr:       true,
			malformJWTStr: false,
		},
		{
			name:          "Malformed JWTStr",
			secret:        root_secret,
			expiresIn:     time.Second * 10,
			wantErr:       true,
			malformJWTStr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtStr, err := MakeJWT(userID, root_secret, tt.expiresIn)
			if err != nil {
				t.Errorf("MakeJWT() error = %v", err)
				return
			}
			if tt.malformJWTStr {
				jwtStr = "malformed.jwt.string"
			}
			userIDFromJWT, err := ValidateJWT(jwtStr, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && userIDFromJWT != userID {
				t.Errorf("userID not match %v != %v", userIDFromJWT, userID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	const headerKey = "Authorization"
	userID := uuid.New()
	root_secret := "MySecret"

	jwtStr, err := MakeJWT(userID, root_secret, time.Hour*10)
	if err != nil {
		t.Errorf("MakeJWT() error = %v", err)
		return
	}

	type headers struct {
		key   string
		value string
	}

	tests := []struct {
		name    string
		headers http.Header
		wantErr bool
	}{
		{
			name: "Valid auth header",
			headers: http.Header{
				headerKey: []string{"Bearer " + jwtStr},
			},
			wantErr: false,
		},
		{
			name:    "Empty auth header",
			headers: http.Header{},
			wantErr: true,
		},
		{
			name: "Invalid auth header length",
			headers: http.Header{
				headerKey: []string{"Bearer " + jwtStr + " asdfasf"},
			},
			wantErr: true,
		},
		{
			name: "Invalid format (excludes 'Bearer ')",
			headers: http.Header{
				headerKey: []string{jwtStr + " slice_length_is_2"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && jwtStr != token {
				t.Errorf("GetBearerToken() token not match %v != %v", jwtStr, token)
			}
		})
	}

}
