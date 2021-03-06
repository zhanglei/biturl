package controllers

import (
	"net/http"
	"strings"

	"github.com/TimothyYe/biturl/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
)

var client *redis.Client

const (
	domain   = "biturl.top"
	url      = "https://biturl.top/"
	visitKey = `visit/%s`
)

//IndexController for URL shorten handling
type IndexController struct {
}

//Response struct for http response
type Response struct {
	Result  bool   `json:"result"`
	Short   string `json:"short"`
	Message string `json:"message"`
}

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

//IndexHandler for rendering the index page
func (c *IndexController) IndexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

//GetShortHandler for getting shorten URL querying result
func (c *IndexController) GetShortHandler(ctx *gin.Context) {
	url := ctx.Param("url")
	longURL := client.Get(url).Val()

	if len(longURL) > 0 {
		if strings.HasPrefix(longURL, "http://") || strings.HasPrefix(longURL, "https://") {
			ctx.Redirect(http.StatusTemporaryRedirect, longURL)
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, "https://"+longURL)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

//ShortURLHandler for shorten long URL
func (c *IndexController) ShortURLHandler(ctx *gin.Context) {
	url := ctx.PostForm("url")
	resp := new(Response)
	inputURL := string(url)

	if !strings.HasPrefix(inputURL, "http") {
		inputURL = "https://" + inputURL
	}

	if inputURL == "" {
		resp.Result = false
		resp.Message = "Please input URL first..."

		ctx.JSON(http.StatusOK, resp)
		return
	}

	if strings.Contains(inputURL, domain) {
		resp.Result = false
		resp.Message = "Cannot shorten it again..."

		ctx.JSON(http.StatusOK, resp)
		return
	}

	urls := utils.ShortenURL(inputURL)
	err := client.Set(urls[0], inputURL, 0).Err()
	if err != nil {
		resp.Result = false
		resp.Message = "Backend service is unavailable!"
	}

	resp.Result = true
	resp.Short = url + urls[0]

	ctx.JSON(http.StatusOK, resp)
}
