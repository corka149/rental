package app

import (
	"log"

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
			c.Status(201)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println(err)
			c.Status(201)
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
			c.Status(201)
			return
		}

		log.Println("User created")
		c.Status(201)
	}
}

func loginForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		templates.Layout(templates.Login()).Render(c.Request.Context(), c.Writer)
	}
}

func login(quries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		email := c.PostForm("email")
		password := c.PostForm("password")

		if email == "" || password == "" {
			log.Println("Email or password is empty", email, password)
			c.Status(200)
			return
		}

		userData, err := quries.GetUserByEmail(c, email)

		if err != nil {
			log.Println(err)
			c.Redirect(200, "/auth/login")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))

		if err != nil {
			log.Println(err)
			c.Redirect(200, "/auth/login")
			return
		}

		session := sessions.Default(c)
		session.Set("user", userData.ID)
		session.Save()

		c.Redirect(302, "/")
	}
}

func logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Status(200)
	}
}
