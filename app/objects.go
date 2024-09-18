package app

import (
	"log"
	"strconv"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
)

func indexObjects(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		objects, err := queries.GetObjects(c.Request.Context())

		if err != nil {
			log.Printf("Error getting objects: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.ObjectIndexes(objects)).Render(c.Request.Context(), c.Writer)
	}
}

func newObjectForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		object := datastore.Object{}
		templates.Layout(user.Name, templates.ObjectForm(object, "new")).Render(c.Request.Context(), c.Writer)
	}
}

func createObject(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")

		_, err := queries.CreateObject(c.Request.Context(), name)
		if err != nil {
			log.Printf("Error creating object: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		c.Redirect(302, "/objects")
	}
}

func updateObjectForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error parsing id: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		object, err := queries.GetObjectById(c.Request.Context(), int32(id))
		if err != nil {
			log.Printf("Error getting object: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.ObjectForm(object, idStr)).Render(c.Request.Context(), c.Writer)
	}
}

func updateObject(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		name := c.PostForm("name")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing id: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		updateObject := datastore.UpdateObjectParams{
			ID:   int32(id),
			Name: name,
		}

		_, err = queries.UpdateObject(c.Request.Context(), updateObject)
		if err != nil {
			log.Printf("Error updating object: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		c.Redirect(302, "/objects")
	}
}

func deleteObject(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error parsing id: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		_, err = queries.DeleteObject(c.Request.Context(), int32(id))
		if err != nil {
			log.Printf("Error deleting object: %v", err)
			c.Redirect(302, "/objects")
			c.Abort()
			return
		}

		c.Redirect(302, "/objects")
	}
}
