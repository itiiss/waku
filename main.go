package main

import (
	"log"
	"net/http"
	"webFramework/waku"
)

func main() {
	r := waku.NewEngine()

	r.Get("/", func(c *waku.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Waku</h1>")
	})
	r.Get("/hello", func(c *waku.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.Get("/hello/:name", func(c *waku.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.Get("/assets/*filepath", func(c *waku.Context) {
		c.JSON(http.StatusOK, waku.H{"filepath": c.Param("filepath")})
	})

	r.Post("/login", func(c *waku.Context) {
		c.JSON(http.StatusOK, waku.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	log.Fatal(r.Run(":9999"))

}
