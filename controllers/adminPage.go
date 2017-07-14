package controllers

import (
	"github.com/akkuman/gooBlog/models"

	"gopkg.in/macaron.v1"
)

func LoginPage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "后台管理登录-" + sitename
	ctx.Data["showError"] = false

	ctx.HTML(200, "login")
}

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

func AdminPage(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isAdd"] = true
		ctx.HTML(200, "admin_add")
	} else {
		ctx.Redirect("/login", 302)
	}
}

func AddArticle(ctx *macaron.Context) {
	if checkLoginStatus(ctx) {
		ctx.Data["pageTitle"] = "后台管理-" + sitename
		ctx.Data["isAdd"] = true
		articleTitle := ctx.Query("articleTitle")
		articleSummary := ctx.Query("articleSummary")
		articleContent := ctx.Query("articleContent")
		ctx.Data["isPostReq"] = true
		if articleTitle != "" {
			var article models.Article
			article.Title = articleTitle
			article.Content = articleContent
			article.Summary = articleSummary
			_ = models.AddArticle(article)
			ctx.Data["titleExist"] = true
			ctx.HTML(200, "admin_add")
		} else {
			ctx.Data["titleExist"] = false
			ctx.HTML(200, "admin_add")
		}
	} else {
		ctx.Redirect("/login", 302)
	}
}

func checkLoginStatus(ctx *macaron.Context) bool {
	username := ctx.GetCookie("username")
	password := ctx.GetCookie("password")
	if username == auth_username && password == auth_password {
		return true
	} else {
		return false
	}
}
