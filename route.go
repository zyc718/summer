package main

import "github.com/gin-gonic/gin"

func route() {
	r := gin.Default()
	apiCommon(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func apiCommon(r *gin.Engine) {
	r.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}