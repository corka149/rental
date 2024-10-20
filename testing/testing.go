// / testing contains test setup and teardown functions.
package rentaltesting

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/corka149/rental/cmd"
	"github.com/corka149/rental/datastore"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func Setup() *cmd.Server {
	gin.SetMode(gin.TestMode)
	ctx := context.TODO()

	s, err := cmd.NewServer(ctx, getenv)

	if err != nil {
		panic(err)
	}

	// Run clean up before running tests
	Teardown(s.Queries)

	return s
}

func Teardown(queries *datastore.Queries) {
	users, _ := queries.GetUsers(context.TODO())

	for _, user := range users {
		queries.DeleteUser(context.TODO(), user.ID)
	}

	rentals, _ := queries.GetRentals(context.TODO())

	for _, rental := range rentals {
		queries.DeleteRental(context.TODO(), rental.ID)
	}

	holidays, _ := queries.GetHolidays(context.TODO())

	for _, holiday := range holidays {
		queries.DeleteHoliday(context.TODO(), holiday.ID)
	}

	objects, _ := queries.GetObjects(context.TODO())

	for _, object := range objects {
		queries.DeleteObject(context.TODO(), object.ID)
	}
}

func getenv(key string) string {
	if key == "DB_URL" {
		return "postgres://myadmin:mypassword@localhost:5432/test_rental_db"
	}

	return ""
}

func CreateUser(queries *datastore.Queries, password string, mod ...func(*datastore.CreateUserParams)) datastore.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	params := &datastore.CreateUserParams{
		Name:     "testUser",
		Email:    "user@test.org",
		Password: string(hashedPassword),
	}

	for _, m := range mod {
		m(params)
	}

	user, err := queries.CreateUser(context.TODO(), *params)

	if err != nil {
		panic(err)
	}

	return user
}

func CreateObject(queries *datastore.Queries, name string) datastore.Object {
	object, err := queries.CreateObject(context.TODO(), name)

	if err != nil {
		panic(err)
	}

	return object
}

func CreateHoliday(queries *datastore.Queries, mod ...func(*datastore.CreateHolidayParams)) datastore.Holiday {
	now := time.Now()

	params := &datastore.CreateHolidayParams{
		Title:     "Holiday",
		Beginning: pgtype.Date{Time: now, Valid: true},
		Ending:    pgtype.Date{Time: now.AddDate(0, 0, 1), Valid: true},
	}

	for _, m := range mod {
		m(params)
	}

	holiday, err := queries.CreateHoliday(context.TODO(), *params)

	if err != nil {
		panic(err)
	}

	return holiday
}

func CreateRental(queries *datastore.Queries, objectid int32, mod ...func(*datastore.CreateRentalParams)) datastore.Rental {
	now := time.Now()

	params := &datastore.CreateRentalParams{
		ObjectID:    objectid,
		Description: pgtype.Text{String: "Rental", Valid: true},
		Beginning:   pgtype.Date{Time: now, Valid: true},
		Ending:      pgtype.Date{Time: now.AddDate(0, 0, 1), Valid: true},
	}

	for _, m := range mod {
		m(params)
	}

	rental, err := queries.CreateRental(context.TODO(), *params)

	if err != nil {
		panic(err)
	}

	return rental
}

func Login(s *cmd.Server, wanted *http.Request, email, password string) {
	w := httptest.NewRecorder()

	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	// Act
	req, err := http.NewRequest("POST", "/auth/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		panic(err)
	}

	s.Router.ServeHTTP(w, req)

	if http.StatusFound != w.Code {
		panic("Login failed")
	}

	session := w.Header().Get("Set-Cookie")

	if !strings.Contains(session, "rental=") {
		panic("No session")
	}

	wanted.Header.Set("Cookie", session)
}
