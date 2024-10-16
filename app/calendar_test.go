package app_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corka149/rental/app"
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
	assert.Equal(t, http.StatusOK, w.Code)

	var entries []app.CalendarEntry

	err := json.Unmarshal(w.Body.Bytes(), &entries)

	assert.Nil(t, err)

	assert.Len(t, entries, 2)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestSearchCalendarEntriesWhenObjectExistsButNoEntriesExists(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()
	now := time.Now()

	name := "Test_Object_1"
	o := rentaltesting.CreateObject(s.Queries, name)

	// Act
	url := fmt.Sprintf("/calendar/search?object=%d&month=%d&year=%d", o.ID, now.Month(), now.Year())

	req, _ := http.NewRequest("GET", url, nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var entries []app.CalendarEntry

	err := json.Unmarshal(w.Body.Bytes(), &entries)

	assert.Nil(t, err)

	assert.Len(t, entries, 0)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
