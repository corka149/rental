package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rentaltesting "github.com/corka149/rental/testing"
	"github.com/stretchr/testify/assert"
)

func TestGetCalendarUI(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	// Act
	req, _ := http.NewRequest("GET", "/calendar", nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
}

func TestSearchCalendarEntriesWhenObjectExists(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()
	now := time.Now()

	name := "Test_Object_1"
	o := rentaltesting.CreateObject(s.Queries, name)
	_ = rentaltesting.CreateRental(s.Queries, o.ID)

	// Act
	url := fmt.Sprintf("/calendar/search?object=%d&month=%d&year=%d", o.ID, now.Month(), now.Year())

	req, _ := http.NewRequest("GET", url, nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
