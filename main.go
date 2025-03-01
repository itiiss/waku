package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"webFramework/waku"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := waku.NewEngine()
	r.Use(waku.Logger())

	// 设置渲染函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载模版
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Alice", Age: 20}
	stu2 := &student{Name: "Bob", Age: 22}

	r.Get("/", func(c *waku.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.Get("/students", func(c *waku.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", waku.H{
			"title":  "waku",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.Get("/date", func(c *waku.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", waku.H{
			"title": "waku",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	log.Fatal(r.Run(":9999"))

}
