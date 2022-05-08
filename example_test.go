package routing_test

import (
	"log"
	"net/http"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/ozzo-routing/v2/access"
	"github.com/zekrotja/ozzo-routing/v2/content"
	"github.com/zekrotja/ozzo-routing/v2/fault"
	"github.com/zekrotja/ozzo-routing/v2/file"
	"github.com/zekrotja/ozzo-routing/v2/slash"
)

func Example() {
	router := routing.New()

	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// serve RESTful APIs
	api := router.Group("/api")
	api.Use(
		// these handlers are shared by the routes in the api group only
		content.TypeNegotiator(content.JSON, content.XML),
	)
	api.Get("/users", func(c *routing.Context) error {
		return c.Write("user list")
	})
	api.Post("/users", func(c *routing.Context) error {
		return c.Write("create a new user")
	})
	api.Put(`/users/<id:\d+>`, func(c *routing.Context) error {
		return c.Write("update user " + c.Param("id"))
	})

	// serve index file
	router.Get("/", file.Content("ui/index.html"))
	// serve files under the "ui" subdirectory
	router.Get("/*", file.Server(file.PathMap{
		"/": "/ui/",
	}))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
