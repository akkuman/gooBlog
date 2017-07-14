package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Article struct {
	gorm.Model
	Title      string `gorm:"size:1<<16-1"`
	Content    string `gorm:"size:1<<16-1"`
	Summary    string `gorm:"size:1<<16-1"`
	Keywords   string
	Url        string
	tag        []Tag
	category   Category
	ReadTime   int
	ViewNum    int
	CommentNum int
}

type Tag struct {
	ID   int
	Name string `gorm:"index"`
}

type Category struct {
	ID   int
	Name string `gorm:"index"`
}

// 创建并初始化数据库
func init() {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	//数据库自动迁移
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Category{})
}

// 添加文章
func AddArticle(updArt Article) uint {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	updArt.ViewNum = 0
	updArt.ReadTime = len(updArt.Content) / 600
	db.Create(&updArt)
	updArt.Url = fmt.Sprintf("/article/%d", updArt.ID)
	db.Model(&updArt).Update("Url", updArt.Url)
	return updArt.ID
}

// 更新文章
func UpdateArticle(updArt Article) {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Model(&updArt).Updates(updArt)
}

// 根据文章id取得该id文章的结构体数据
func GetArticle(id uint) Article {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var article Article
	db.First(&article, id)
	db.Model(&article).Update("ViewNum", article.ViewNum+1)
	return article
}

func ListArticle(pagenum int, pageOffset int) []Article {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var articles []Article
	if pagenum > 0 {
		db.Order("id desc").Offset((pagenum - 1) * pageOffset).Limit(pageOffset).Find(&articles)
	} else {
		db.Order("id desc").Offset(-1).Limit(pageOffset).Find(&articles)
	}
	return articles
}

// 取得所有的文章数
func GetArticleNum() interface{} {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var count interface{}
	db.Model(&Article{}).Count(&count)
	return count
}

// 传入本篇文章id，取得前一篇文章的结构体数据
func GetPrevArticle(id uint) Article {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var articles []Article
	db.Where(fmt.Sprintf("id < %d", id)).Order("id desc").Limit(1).Find(&articles)
	if len(articles) != 0 {
		return articles[0]
	}
	return Article{}
}

// 传入本篇文章id，取得后一篇文章的结构体数据
func GetNextArticle(id uint) Article {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var articles []Article
	db.Where(fmt.Sprintf("id > %d", id)).Limit(1).Find(&articles)
	if len(articles) != 0 {
		return articles[0]
	}
	return Article{}
}
