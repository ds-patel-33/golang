package main

import "github.com/gin-gonic/gin"

var router *gin.Engine

func initializeRoutes() {

	router.GET("/", showIndexPage)

	router.GET("/article/view/:article_id", getArticle)

}
