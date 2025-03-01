package main

import (
	"log"
	"net/http"
	"webFramework/waku"
)

func main() {
	r := waku.NewEngine()

	r.Get("/index", func(c *waku.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.NewGroup("/v1")
	{
		v1.Get("/", func(c *waku.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Waku</h1>")
		})

		v1.Get("/hello", func(c *waku.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.NewGroup("/v2")
	{
		v2.Get("/hello/:name", func(c *waku.Context) {

			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.Post("/login", func(c *waku.Context) {
			c.JSON(http.StatusOK, waku.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	log.Fatal(r.Run(":9999"))

}
