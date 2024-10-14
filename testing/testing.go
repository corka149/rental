// / testing contains test setup and teardown functions.
package rentaltesting

import (
	"context"
	"time"

	"github.com/corka149/rental/cmd"
	"github.com/corka149/rental/datastore"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func Setup() *cmd.Server {
	ctx := context.TODO()

	s, err := cmd.NewServer(ctx, getenv)

	if err != nil {
		panic(err)
	}

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
	params := &datastore.CreateHolidayParams{
		Title: "Holiday",
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
