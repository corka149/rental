package app_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/corka149/rental/datastore"
	rentaltesting "github.com/corka149/rental/testing"
	"github.com/jackc/pgx/v5/pgtype"
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

func TestNewHolidayForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/holidays/new", nil)
	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestCreateHoliday(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	beginning := "2024-10-20"
	ending := "2024-10-30"

	formData := url.Values{}
	formData.Set("beginning", beginning)
	formData.Set("ending", ending)
	formData.Set("title", "Test Holiday")

	// Act
	req, err := http.NewRequest("POST", "/holidays/new", strings.NewReader(formData.Encode()))

	assert.Nil(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	holidays, err := s.Queries.GetHolidays(context.Background())

	assert.Nil(t, err)
	assert.Len(t, holidays, 1)

	assert.Equal(t, "Test Holiday", holidays[0].Title)
	assert.Equal(t, beginning, holidays[0].Beginning.Time.Format("2006-01-02"))
	assert.Equal(t, ending, holidays[0].Ending.Time.Format("2006-01-02"))

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestCreateHolidayFailsWhenItOverlaps(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	beginning := "2024-10-20"
	ending := "2024-10-30"

	existing := rentaltesting.CreateHoliday(s.Queries, func(params *datastore.CreateHolidayParams) {
		params.Beginning = pgtype.Date{Time: time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC), Valid: true}
		params.Ending = pgtype.Date{Time: time.Date(2024, 10, 27, 0, 0, 0, 0, time.UTC), Valid: true}
	})

	formData := url.Values{}
	formData.Set("beginning", beginning)
	formData.Set("ending", ending)
	formData.Set("title", "Test Holiday")

	// Act
	req, err := http.NewRequest("POST", "/holidays/new", strings.NewReader(formData.Encode()))

	assert.Nil(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "There is a conflict with a holiday.")

	holidays, err := s.Queries.GetHolidays(context.Background())

	assert.Nil(t, err)
	assert.Len(t, holidays, 1)
	assert.Equal(t, existing.ID, holidays[0].ID)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateHolidayForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	holiday := rentaltesting.CreateHoliday(s.Queries)

	// Act
	url := fmt.Sprintf("/holidays/%d", holiday.ID)

	req, _ := http.NewRequest("GET", url, nil)

	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateHoliday(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	holiday := rentaltesting.CreateHoliday(s.Queries)

	beginning := "2024-10-20"
	ending := "2024-10-30"

	formData := url.Values{}
	formData.Set("beginning", beginning)
	formData.Set("ending", ending)
	formData.Set("title", "Test Holiday")

	// Act
	url := fmt.Sprintf("/holidays/%d", holiday.ID)

	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))

	assert.Nil(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	holidays, err := s.Queries.GetHolidays(context.Background())

	assert.Nil(t, err)
	assert.Len(t, holidays, 1)

	assert.Equal(t, "Test Holiday", holidays[0].Title)
	assert.Equal(t, beginning, holidays[0].Beginning.Time.Format("2006-01-02"))
	assert.Equal(t, ending, holidays[0].Ending.Time.Format("2006-01-02"))

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestDeleteHoliday(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	holiday := rentaltesting.CreateHoliday(s.Queries)

	// Act
	url := fmt.Sprintf("/holidays/%d/delete", holiday.ID)

	req, _ := http.NewRequest("POST", url, nil)

	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	holidays, err := s.Queries.GetHolidays(context.Background())

	assert.Nil(t, err)
	assert.Len(t, holidays, 0)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
