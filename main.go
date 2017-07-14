package main

import (
	"html/template"
	"log"
	"net/http"

	. "github.com/akkuman/gooBlog/controllers"

	"github.com/goless/config"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.Classic()
	//从配置文件设定模板目录;设定静态文件目录
	configure := config.New("config.json")
	templateDir := "templates/" + configure.Get("theme").(string)
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: templateDir,
		Funcs: []template.FuncMap{map[string]interface{}{
			"RawHTML": templateFunc_RawHTML,
		}},
	}))
	m.Use(macaron.Static(templateDir))
	//路由
	m.Get("/", BasicRenderData, HomePage)
	m.Get("/archive", BasicRenderData, ArchivePage)
	m.Get("/about", BasicRenderData, AboutPage)
	m.Get("/page/:pagenum(^[1-9]\\d*$)", BasicRenderData, HomePage)
	m.Get("/admin", BasicRenderData, AdminPage)
	m.Combo("/login").Get(BasicRenderData, LoginPage).Post(BasicRenderData, CheckLogin)
	m.Get("/admin/article", BasicRenderData, AdminPage)
	m.Combo("/admin/article/add").Get(BasicRenderData, AdminPage).Post(BasicRenderData, PostAddArticle)
	m.Get("/admin/manager", BasicRenderData, ManagerPage)
	m.Post("/admin/category/add", BasicRenderData, PostAddCategory)
	m.Post("/admin/tag/add", BasicRenderData, PostAddTag)
	m.Get("/article/:id([0-9]+)", BasicRenderData, ArticlePage)

	m.NotFound(BasicRenderData, NotFoundPage)

	log.Println("Server is running...")
	log.Println(http.ListenAndServe(configure.Get("ListenAndServe").(string), m))
}

func templateFunc_RawHTML(parsedHTML string) interface{} {
	return template.HTML(parsedHTML)
}
