package app_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	rentaltesting "github.com/corka149/rental/testing"
	"github.com/stretchr/testify/assert"
)

func TestGetObjectsIndexWithoutData(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/objects", nil)
	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestGetObjectsIndexWithData(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)
	o := rentaltesting.CreateObject(s.Queries, "test_object_")

	// Act
	req, err := http.NewRequest("GET", "/objects", nil)

	assert.NoError(t, err)

	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), o.Name)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestNewObjectForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)

	// Act
	req, _ := http.NewRequest("GET", "/objects/new", nil)
	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestCreateObject(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)

	name := "test_object_"

	formData := url.Values{}
	formData.Set("name", name)

	// Act
	req, err := http.NewRequest("POST", "/objects/new", strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.NoError(t, err)

	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	objects, err := s.Queries.GetObjects(context.Background())

	assert.NoError(t, err)

	assert.Len(t, objects, 1)
	assert.Equal(t, name, objects[0].Name)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateObjectForm(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)
	o := rentaltesting.CreateObject(s.Queries, "test_object_")

	// Act
	url := fmt.Sprintf("/objects/%d", o.ID)

	req, _ := http.NewRequest("GET", url, nil)
	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), o.Name)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestUpdateObject(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)
	o := rentaltesting.CreateObject(s.Queries, "test_object_")

	name := "test_object_updated"

	formData := url.Values{}
	formData.Set("name", name)

	// Act
	url := fmt.Sprintf("/objects/%d", o.ID)

	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.Nil(t, err)

	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	objects, err := s.Queries.GetObjects(context.Background())

	assert.Nil(t, err)
	assert.Len(t, objects, 1)

	assert.Equal(t, name, objects[0].Name)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}

func TestDeleteObject(t *testing.T) {
	// Arrange
	s := rentaltesting.Setup()
	w := httptest.NewRecorder()

	password := "password"
	u := rentaltesting.CreateUser(s.Queries, password)
	o := rentaltesting.CreateObject(s.Queries, "test_object_")

	// Act
	url := fmt.Sprintf("/objects/%d/delete", o.ID)

	req, _ := http.NewRequest("POST", url, nil)

	rentaltesting.Login(s, req, u.Email, password)

	s.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusFound, w.Code)

	objects, err := s.Queries.GetObjects(context.Background())

	assert.Nil(t, err)
	assert.Len(t, objects, 0)

	// Teardown
	rentaltesting.Teardown(s.Queries)
}
