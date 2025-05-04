package main

import (
	"gee/gee"
	"log"
	"net/http"
)

func forAbort() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("Aborting....")
		c.Abort()
	}
}

func forNext1() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("Next1 before....")
		c.Next()
		log.Printf("Next1 after....")
	}
}

func forNext2() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("Next2 before....")
		c.Next()
		log.Printf("Next2 after....")
	}
}

func main() {
	r := gee.Default()

	r.GET("/index", func(c *gee.Context) {
		log.Println("do /index")
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	r.GET("/recovery", func(c *gee.Context) {
		log.Println("do /recovery")
		panic("crash!!!!!")
		c.Text(http.StatusOK, "recover failed....")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			log.Println("do /")
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			log.Println("do /hello")
			// expect /hello?name=geektutu
			c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

		a := v1.Group("/abort")
		a.Use(forNext1(), forAbort(), forNext2())
		a.GET("/abortNow", func(c *gee.Context) {
			log.Println("do /abortNow")
			c.Text(http.StatusOK, "fail to abort....")
		})
	}
	v2 := r.Group("/v2")
	v2.Use(forNext1(), forNext2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			log.Println("do /hello/:name")
			// expect /hello/geektutu
			c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			log.Println("do /login")
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")
}
