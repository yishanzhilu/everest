package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func foo(c *gin.Context) {

}

// Roadmap .
type Roadmap struct {
	Steps []interface{} `json:"steps"`
}

// StudyStep .
type StudyStep interface {
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()

	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	content := []byte(`{"steps":[{"title":"freeCodeCamp HTML 基础","description":"","type":"online-course","options":[{"locales":"en","url":"https://www.freecodecamp.org/learn/responsive-web-design/basic-html-and-html5/"},{"locales":"zh-CN","url":"https://learn.freecodecamp.one/responsive-web-design/basic-html-and-html5"}]},{"title":"CSS","description":"","type":"parallel-topic","steps":[{"title":"freeCodeCamp CSS 基础","description":"","type":"web-page","options":[{"locales":"en","url":"https://www.freecodecamp.org/learn/responsive-web-design/basic-css/"},{"locales":"zh-CN","url":"https://learn.freecodecamp.one/responsive-web-design/basic-css"}]},{"title":"CSS Flex 布局","description":"","type":"web-page","options":[{"locales":"en","url":"https://css-tricks.com/snippets/css/a-guide-to-flexbox/"},{"locales":"zh-CN","url":"http://www.ruanyifeng.com/blog/2015/07/flex-grammar.html"}]}]}]}`)
	r := &Roadmap{}
	if err := json.Unmarshal(content, r); err != nil {
		panic(err)
	}
	fmt.Println("test ", r)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, r)
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
