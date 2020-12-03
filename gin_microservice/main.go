// main.go

package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	router.Run()

}

func initializeRoutes() {

	router.GET("/", showIndexPage)
	router.GET("/about-us", aboutUsPage)
	router.GET("/contact-us", contactUsPage)

	router.GET("/article/view/:article_id", getArticle)

}

func showIndexPage(c *gin.Context) {
	articles := getAllArticles()

	c.HTML(

		http.StatusOK,

		"index.html",

		gin.H{
			"title":   "Home Page",
			"payload": articles,
		},
	)

}

func contactUsPage(c *gin.Context) {
	articles := getAllArticles()

	c.HTML(

		http.StatusOK,

		"contact_us.html",

		gin.H{
			"title":   "Home Page",
			"payload": articles,
		},
	)

}

func aboutUsPage(c *gin.Context) {
	articles := getAllArticles()

	c.HTML(

		http.StatusOK,

		"about_us.html",

		gin.H{
			"title":   "Home Page",
			"payload": articles,
		},
	)

}

func getArticle(c *gin.Context) {

	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {

		if article, err := getArticleByID(articleID); err == nil {

			c.HTML(

				http.StatusOK,

				"article.html",

				gin.H{
					"title":   article.Title,
					"payload": article,
				},
			)

		} else {

			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {

		c.AbortWithStatus(http.StatusNotFound)
	}
}

type article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articleList = []article{
	article{ID: 1, Title: "Article 1", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
	article{ID: 2, Title: "Article 2", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
	article{ID: 3, Title: "Article 3", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
	article{ID: 4, Title: "Article 4", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
	article{ID: 5, Title: "Article 5", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
	article{ID: 6, Title: "Article 6", Content: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."},
}

func getAllArticles() []article {
	return articleList
}

func getArticleByID(id int) (*article, error) {
	for _, a := range articleList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Article not found")
}
