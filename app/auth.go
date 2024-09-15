package app

import (
	"log"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/dto"
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
