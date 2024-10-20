package app

import (
	"log"
	"net/http"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/dto"
	"github.com/corka149/rental/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func register(quries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := dto.CreateUser{}

		err := c.BindJSON(&user)

		if err != nil {
			log.Println(err)
			c.Status(http.StatusCreated)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println(err)
			c.Status(http.StatusCreated)
			return
		}

		createParams := datastore.CreateUserParams{
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hashedPassword),
		}

		_, err = quries.CreateUser(c, createParams)

		if err != nil {
			log.Println(err)
			c.Status(http.StatusCreated)
			return
		}

		log.Println("User created")
		c.Status(http.StatusCreated)
	}
}

func loginForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)

		templates.Layout(user.Name, templates.Login()).Render(c.Request.Context(), c.Writer)
	}
}

func login(quries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		email := c.PostForm("email")
		password := c.PostForm("password")

		if email == "" || password == "" {
			log.Println("Email or password is empty", email, password)
			c.Status(http.StatusOK)
			return
		}

		userData, err := quries.GetUserByEmail(c, email)

		if err != nil {
			log.Printf("Error getting user by email: %v", err)
			c.Redirect(http.StatusOK, "/auth/login")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))

		if err != nil {
			log.Printf("Password not match: %v", err)
			c.Status(http.StatusUnauthorized)
			templates.Layout("", templates.Login()).Render(c.Request.Context(), c.Writer)
			return
		}

		session := sessions.Default(c)
		session.Set("user", userData.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}

func logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	}
}

func getUserFromSession(c *gin.Context, queries *datastore.Queries) datastore.User {
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int32)

	if !ok {
		return datastore.User{}
	}

	user, err := queries.GetUserById(c, userID)

	if err != nil {
		log.Printf("Error getting user from session: %v", err)
		return datastore.User{}
	}

	return user
}
