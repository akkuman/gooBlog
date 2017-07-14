package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akkuman/gooBlog/models"

	"github.com/goless/config"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/macaron.v1"
)

func BasicRenderData(ctx *macaron.Context) {
	ctx.Data["siteName"] = sitename
	ctx.Data["quoteSlice"] = quote_slice
	ctx.Data["linkGithub"] = link_github
	ctx.Data["linkTwitter"] = link_twitter
	ctx.Data["author"] = author
	ctx.Data["copyright"] = copyright
}

func HomePage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "首页-" + sitename
	ctx.Data["isHome"] = true
	var pagenum int
	if pagenumParam := ctx.Params(":pagenum"); pagenumParam == "" {
		pagenum = 1
	} else {
		pagenum, _ = strconv.Atoi(pagenumParam)
	}
	pageOffset, _ := strconv.Atoi(getConfigString("page_offset"))
	ctx.Data["articles"] = models.ListArticle(pagenum, pageOffset)
	// 判断前一页是否还有数据
	if pagenum-1 > 0 {
		ctx.Data["hasNewerArticles"] = true
		ctx.Data["newerArticlesLink"] = fmt.Sprintf("/page/%d", pagenum-1)
	}
	// 判断后一页是否还有数据
	if len(models.ListArticle(pagenum+1, pageOffset)) != 0 {
		ctx.Data["hasOlderArticles"] = true
		ctx.Data["olderArticlesLink"] = fmt.Sprintf("/page/%d", pagenum+1)
	}

	ctx.HTML(200, "home")
}

func ArchivePage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "归档-" + sitename
	ctx.Data["isArchive"] = true
	ctx.Data["blogsNum"] = models.GetArticleNum()
	ctx.Data["todayDate"] = time.Now().Format("2006-01-02")
	ctx.Data["articles"] = models.ListArticle(-1, -1)

	ctx.HTML(200, "archive")
}

func ArticlePage(ctx *macaron.Context) {
	id, _ := strconv.Atoi(ctx.Params(":id"))
	article := models.GetArticle(uint(id))
	ctx.Data["pageTitle"] = article.Title + "-" + sitename
	article.Content = StringMarkdownToHTML(article.Content)
	PrevArticle := models.GetPrevArticle(uint(id))
	NextArticle := models.GetNextArticle(uint(id))

	if PrevArticle.Url != "" {
		ctx.Data["hasPrevArticle"] = true
	}
	if NextArticle.Url != "" {
		ctx.Data["hasNextArticle"] = true
	}
	ctx.Data["isArticle"] = true
	ctx.Data["article"] = article
	ctx.Data["prevArticle"] = PrevArticle
	ctx.Data["nextArticle"] = NextArticle
	ctx.HTML(200, "article")
}

func AboutPage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "关于我-" + sitename
	ctx.Data["aboutSlice"] = about_slice
	ctx.Data["blogsNum"] = models.GetArticleNum()
	ctx.Data["isAbout"] = true

	ctx.HTML(200, "about")
}

func NotFoundPage(ctx *macaron.Context) {
	ctx.Data["pageTitle"] = "页面没找到-" + sitename

	ctx.HTML(200, "404")

}

func getConfigString(item string) string {
	c := config.New("config.json")
	return c.Get(item).(string)
}

func StringMarkdownToHTML(textData string) string {
	byteData := blackfriday.MarkdownCommon([]byte(textData))
	HTMLData := string(byteData[:])
	HTMLData = bluemonday.UGCPolicy().Sanitize(HTMLData)
	return HTMLData
}
