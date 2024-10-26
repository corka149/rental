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

func TestIndexRentals(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/rentals", nil)
	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestIndexRentalsRequiresAuth(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	// Act
	req, _ := http.NewRequest("GET", "/rentals", nil)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	assert.Contains(t, w.Header().Get("Location"), "/auth/login")

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestNewRentalForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/rentals/new", nil)
	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRental(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	beginning := "2024-10-20"
	ending := "2024-10-30"

	formData := url.Values{}
	formData.Set("beginning", beginning)
	formData.Set("ending", ending)
	formData.Set("object", fmt.Sprintf("%d", object.ID))
	formData.Set("description", "school")

	// Existing of other rental should not matter
	otherObject := rentaltesting.CreateObject(s.Queries, "other_object_")
	rentaltesting.CreateRental(s.Queries, otherObject.ID, func(params *datastore.CreateRentalParams) {
		params.Beginning = pgtype.Date{Time: time.Date(2024, 10, 19, 0, 0, 0, 0, time.UTC), Valid: true}
		params.Ending = pgtype.Date{Time: time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC), Valid: true}
	})

	// Act
	req, err := http.NewRequest("POST", "/rentals/new", strings.NewReader(formData.Encode()))

	assert.NoError(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	assert.Contains(t, w.Header().Get("Location"), "/rentals")

	rentals, err := s.Queries.GetRentals(context.Background())

	assert.NoError(t, err)

	assert.Len(t, rentals, 2)

	var newRental datastore.Rental

	for _, r := range rentals {
		if r.ObjectID == object.ID {
			newRental = r
			break
		}
	}

	assert.Equal(t, object.ID, newRental.ObjectID)
	assert.Equal(t, "school", newRental.Description.String)
	assert.Equal(t, beginning, newRental.Beginning.Time.Format("2006-01-02"))
	assert.Equal(t, ending, newRental.Ending.Time.Format("2006-01-02"))

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestCreateRentalFailsWhenOverlapsWithRental(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	existing := rentaltesting.CreateRental(s.Queries, object.ID)

	formData := url.Values{}
	formData.Set("beginning", existing.Beginning.Time.Format("2006-01-02"))
	formData.Set("ending", existing.Ending.Time.Format("2006-01-02"))
	formData.Set("object", fmt.Sprintf("%d", object.ID))
	formData.Set("description", "school")

	// Act
	req, err := http.NewRequest("POST", "/rentals/new", strings.NewReader(formData.Encode()))

	assert.NoError(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "There is a conflict with a rental.")

	rentals, err := s.Queries.GetRentals(context.Background())

	assert.NoError(t, err)

	assert.Len(t, rentals, 1)

	assert.Equal(t, existing.ID, rentals[0].ID)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestCreateRentalFailsWhenOverlapsWithHoliday(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	holiday := rentaltesting.CreateHoliday(s.Queries)

	formData := url.Values{}
	formData.Set("beginning", holiday.Beginning.Time.Format("2006-01-02"))
	formData.Set("ending", holiday.Ending.Time.Format("2006-01-02"))
	formData.Set("object", fmt.Sprintf("%d", object.ID))
	formData.Set("description", "school")

	// Act
	req, err := http.NewRequest("POST", "/rentals/new", strings.NewReader(formData.Encode()))

	assert.NoError(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "There is a conflict with a holiday.")

	rentals, err := s.Queries.GetRentals(context.Background())

	assert.NoError(t, err)

	assert.Len(t, rentals, 0)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateRentalForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	rental := rentaltesting.CreateRental(s.Queries, object.ID)

	// Act
	url := fmt.Sprintf("/rentals/%d", rental.ID)

	req, _ := http.NewRequest("GET", url, nil)

	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateRental(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	rental := rentaltesting.CreateRental(s.Queries, object.ID)

	beginning := "2030-10-20"
	ending := "2030-10-30"

	formData := url.Values{}
	formData.Set("beginning", beginning)
	formData.Set("ending", ending)
	formData.Set("object", fmt.Sprintf("%d", object.ID))
	formData.Set("description", "school")

	// Act
	url := fmt.Sprintf("/rentals/%d", rental.ID)

	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))

	assert.NoError(t, err)

	rentaltesting.Login(s, req, user.Email, password)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	assert.Contains(t, w.Header().Get("Location"), "/rentals")

	rentals, err := s.Queries.GetRentals(context.Background())

	assert.NoError(t, err)

	assert.Len(t, rentals, 1)

	assert.Equal(t, object.ID, rentals[0].ObjectID)
	assert.Equal(t, "school", rentals[0].Description.String)
	assert.Equal(t, beginning, rentals[0].Beginning.Time.Format("2006-01-02"))
	assert.Equal(t, ending, rentals[0].Ending.Time.Format("2006-01-02"))

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestDeleteRental(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	user := rentaltesting.CreateUser(s.Queries, password)

	object := rentaltesting.CreateObject(s.Queries, "test_object_")

	rental := rentaltesting.CreateRental(s.Queries, object.ID)

	// Act
	url := fmt.Sprintf("/rentals/%d/delete", rental.ID)

	req, _ := http.NewRequest("POST", url, nil)

	rentaltesting.Login(s, req, user.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	rentals, err := s.Queries.GetRentals(context.Background())

	assert.NoError(t, err)

	assert.Len(t, rentals, 0)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
