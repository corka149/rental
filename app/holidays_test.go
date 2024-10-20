package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rentaltesting "github.com/corka149/rental/testing"
	"github.com/stretchr/testify/assert"
)

func TestIndexHolidays(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/holidays", nil)
	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestIndexHolidaysRequiresAuth(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	// Act
	req, _ := http.NewRequest("GET", "/holidays", nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	assert.Contains(t, w.Header().Get("Location"), "/auth/login")

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
