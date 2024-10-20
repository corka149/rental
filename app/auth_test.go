package app_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	rentaltesting "github.com/corka149/rental/testing"
	"github.com/stretchr/testify/assert"
)

func TestShowLoginForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	// Act
	req, _ := http.NewRequest("GET", "/auth/login", nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
}

func TestLogin(t *testing.T) {
	// Arrange
	password := "password"

	s := rentaltesting.Setup()
	w := httptest.NewRecorder()
	u := rentaltesting.CreateUser(s.Queries, password)

	formData := url.Values{}
	formData.Set("email", u.Email)
	formData.Set("password", password)

	// Act
	req, err := http.NewRequest("POST", "/auth/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.Nil(t, err)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	session := w.Header().Get("Set-Cookie")

	assert.Contains(t, session, "rental=")

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestFailsOnWrongPassword(t *testing.T) {
	// Arrange
	password := "password"

	s := rentaltesting.Setup()
	w := httptest.NewRecorder()
	u := rentaltesting.CreateUser(s.Queries, password)

	formData := url.Values{}
	formData.Set("email", u.Email)
	formData.Set("password", "nope")

	// Act
	req, err := http.NewRequest("POST", "/auth/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.Nil(t, err)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	session := w.Header().Get("Set-Cookie")

	assert.NotContains(t, session, "rental=")

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
