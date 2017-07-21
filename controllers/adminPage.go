package controllers

import (
	"time"

	"github.com/akkuman/gooBlog/models"

	"gopkg.in/macaron.v1"
)

// GET登录页面
func LoginPage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "后台管理登录-" + sitename
	ctx.Data["showError"] = false

	ctx.HTML(200, "login")
}

// POST登录数据并检查
func CheckLogin(ctx *macaron.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if username == auth_username && password == auth_password {
		ctx.SetCookie("username", username)
		ctx.SetCookie("password", password)
		ctx.Redirect("/admin", 302)
	} else {
		ctx.Data["showError"] = true
		ctx.HTML(200, "login")
	}
}

// GET检查cookie进入后台
func AdminPage(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isAdd"] = true
		ctx.Data["categories"] = models.GetCategory(uint(0))
		ctx.Data["tags"] = models.GetTags(uint(0))
		ctx.HTML(200, "admin_article")
	} else {
		ctx.Redirect("/login", 302)
	}
}

// POST发表文章
func PostAddArticle(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isAdd"] = true
		ctx.Data["categories"] = models.GetCategory(uint(0))
		ctx.Data["tags"] = models.GetTags(uint(0))
		articleTitle := ctx.Query("articleTitle")
		articleSummary := ctx.Query("articleSummary")
		articleContent := ctx.Query("articleContent")
		articleCateName := ctx.Query("articleCate")
		articleTagsName := ctx.QueryStrings("articleTags")
		ctx.Data["isPostReq"] = true
		if articleTitle != "" {
			ctx.Data["titleExist"] = true
			var article models.Article
			article.Title = articleTitle
			article.Content = articleContent
			article.Summary = articleSummary
			article.Category = models.Category{Name: articleCateName}
			// 根据POST值构建Tag结构列表
			var tags []models.Tag
			for _, v := range articleTagsName {
				tags = append(tags, models.Tag{Name: v})
			}
			article.Tags = tags
			_ = models.AddArticle(article)
		}
		ctx.HTML(200, "admin_article")
	} else {
		ctx.Redirect("/login", 302)
	}
}

// 检查cookie判断是否登录
func checkLoginStatus(ctx *macaron.Context) bool {
	username := ctx.GetCookie("username")
	password := ctx.GetCookie("password")
	if username == auth_username && password == auth_password {
		return true
	} else {
		return false
	}
}

// GET管理页面
func ManagerPage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "后台管理-" + sitename
	ctx.Data["isManager"] = true
	ctx.Data["articles"] = models.ListArticle(-1, -1)

	ctx.HTML(200, "admin_manager")
}

// POST添加分类
func PostAddCategory(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isManager"] = true
		ctx.Data["isPostReq"] = true
		categoryName := ctx.Query("category_name")
		if categoryName != "" {
			ctx.Data["titleExist"] = true
			var category models.Category
			category.Name = categoryName
			_ = models.AddCategory(category)
		}
		ctx.HTML(200, "admin_manager")
	} else {
		ctx.Redirect("/login", 302)
	}
}

// POST添加标签
func PostAddTag(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isManager"] = true
		ctx.Data["isPostReq"] = true
		tagName := ctx.Query("tag_name")
		if tagName != "" {
			ctx.Data["titleExist"] = true
			var tag models.Tag
			tag.Name = tagName
			_ = models.AddTag(tag)
		}
		ctx.HTML(200, "admin_manager")
	} else {
		ctx.Redirect("/login", 302)
	}
}

// 文章修改页面
func EditArticle(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "编辑文章-" + sitename
		ctx.Data["isAdd"] = true
		id := ctx.ParamsInt("id")
		article := models.GetArticle(uint(id))
		ctx.Data["article"] = article
		ctx.HTML(200, "edit_article")
	} else {
		ctx.Redirect("/login", 302)
	}
}

// POST修改文章
func PostEditArticle(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		var article models.Article
		article.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", ctx.Query("articleCreatedAt"))
		article.Title = ctx.Query("articleTitle")
		article.Summary = ctx.Query("articleSummary")
		article.Content = ctx.Query("articleContent")
		if article.Title != "" {
			models.UpdateArticle(article)
		}
		ctx.Redirect("/admin/manager", 302)
	} else {
		ctx.Redirect("/login", 302)
	}
}

func DelArticle(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		id := ctx.ParamsInt("id")
		models.DelArticle(uint(id))
		ctx.Redirect("/admin/manager", 302)
	} else {
		ctx.Redirect("/login", 302)
	}
}
