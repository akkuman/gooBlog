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
	Tags       []Tag    `gorm:"ForeignKey:ArticleID"`
	Category   Category `gorm:"ForeignKey:ArticleID"`
	ReadTime   int
	ViewNum    int
	CommentNum int
}

type Tag struct {
	ID        uint
	ArticleID uint
	Name      string
}

type Category struct {
	ID        uint
	ArticleID uint
	Name      string
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

// 根据id删除指定id文章
func DelArticle(id uint) {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var updArt Article
	updArt.ID = id
	db.Delete(&updArt)
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
	article.Tags = GetTags(id)
	article.Category = GetCategory(id)[0]
	return article
}

// 遍历返回指定页数的文章
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
	db.Where("id < ?", id).Order("id desc").Limit(1).Find(&articles)
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
	db.Where("id > ?", id).Limit(1).Find(&articles)
	if len(articles) != 0 {
		return articles[0]
	}
	return Article{}
}

// 添加分类
func AddCategory(cate Category) uint {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Create(&cate)
	return cate.ID
}

// 取得指定ArticleID分类的数据列表,参数为0取出所有
func GetCategory(articleID uint) []Category {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var categories []Category
	db.Where("article_id = ?", articleID).Find(&categories)
	return categories
}

// 添加标签
func AddTag(tag Tag) uint {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Create(&tag)
	return tag.ID
}

// 取得指定ArticleID标签的数据列表,参数为0取出所有
func GetTags(articleID uint) []Tag {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var tags []Tag
	db.Where("article_id = ?", articleID).Find(&tags)
	return tags
}
