package main

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type json []byte

func (j json) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}
func (j *json) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}
func (j json) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}
func (j *json) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
func (j json) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}
func (j json) Equals(j1 json) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

// Roadmap .
type Roadmap struct {
	gorm.Model
	Title   string `gorm:"type:varchar(80)"`
	Content json   `sql:"type:json" json:"object,omitempty"`
}

func main() {
	db, err := gorm.Open("mysql", "qy:871127@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&Roadmap{})
	// Read
	var r Roadmap
	r.Title = "前端工程师"
	r.Content = []byte(`{"steps":[{"title":"freeCodeCamp HTML 基础","description":"","type":"online-course","options":[{"locales":"en","url":"https://www.freecodecamp.org/learn/responsive-web-design/basic-html-and-html5/"},{"locales":"zh-CN","url":"https://learn.freecodecamp.one/responsive-web-design/basic-html-and-html5"}]},{"title":"CSS","description":"","type":"parallel-topic","steps":[{"title":"freeCodeCamp CSS 基础","description":"","type":"web-page","options":[{"locales":"en","url":"https://www.freecodecamp.org/learn/responsive-web-design/basic-css/"},{"locales":"zh-CN","url":"https://learn.freecodecamp.one/responsive-web-design/basic-css"}]},{"title":"CSS Flex 布局","description":"","type":"web-page","options":[{"locales":"en","url":"https://css-tricks.com/snippets/css/a-guide-to-flexbox/"},{"locales":"zh-CN","url":"http://www.ruanyifeng.com/blog/2015/07/flex-grammar.html"}]}]}]}`)
	db.Save(&r)
}
