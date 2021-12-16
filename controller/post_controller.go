package controller

import (
	"net/http"
	"strconv"
	"time"

	service "personal-blog/service"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title      string    `json:"Title"`
	Body       string    `json:"Body"`
	CreateDate time.Time `json:"CreateDate"`
	Id         string    `json:"Id"`
	CategoryId int       `json:"CategoryID"`
}

// PostExample godoc
// @Summary      post request example
// @Description  post request example
// @Accept       json
// @Produce      plain
// @Param        message  body      bool  true  "Account Info"
// @Success      200      {string}  string         "success"
// @Failure      500      {string}  string         "fail"
// @Router       /posts [post]
func CreatePost(c *gin.Context) {
	var p Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	result, _ := service.Add(&service.PostDto{CategoryId: p.CategoryId, Title: p.Title, Body: p.Body, CreateDate: time.Now()})
	if result {
		c.IndentedJSON(http.StatusOK, result)
		return
	}
}

func UpdatePost(c *gin.Context) {

}

func GetPostById(c *gin.Context) {
	id := c.Param("id")

	post := service.FindById(id)

	if (service.PostDto{}) != post {
		c.IndentedJSON(http.StatusOK, post)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}

func GetPosts(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	if page <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page cannot be lower than 0"})
		return
	}
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
	if size <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "size cannot be lower than 0"})
		return
	}

	posts := service.FindByPagination(page, size)

	if posts != nil {
		c.IndentedJSON(http.StatusOK, posts)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}

func GetPostsByCategoryId(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("categoryid"), 10, 64)

	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	if page <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page cannot be lower than 0"})
		return
	}

	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
	if size <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "size cannot be lower than 0"})
		return
	}

	posts := service.FindByCategoryPagination(id, page, size)
	if posts != nil {
		c.IndentedJSON(http.StatusOK, posts)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "posts not found"})
}
