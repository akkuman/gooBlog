package controllers

import (
	"github.com/goless/config"
)

var configure = config.New("config.json")
var auth_username = configure.Get("auth.username").(string)
var auth_password = configure.Get("auth.password").(string)
var sitename = configure.Get("sitename").(string)
var quote_slice = configure.Get("quote").([]interface{})
var about_slice = configure.Get("about").([]interface{})
var author = configure.Get("author").(string)
var copyright = configure.Get("copyright").(string)
var link_github = "https://github.com/" + configure.Get("social.Github").(string)
var link_twitter = "https://twitter.com/" + configure.Get("social.Twitter").(string)

//var link_weibo = "http://weibo.com/" + configure.Get("social.Weibo").(string)

func configSlice(interfaceSlice []interface{}) []string {
	var stringSlice []string
	for _, v := range interfaceSlice {
		stringSlice = append(stringSlice, v.(string))
	}
	return stringSlice
}
