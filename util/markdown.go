package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"
)

func FetchRedisData() ([]Movie, error) {
	skey := time.Now().Format("movie-2006-01-02")
	ids, err := redisClient.SMembers(skey).Result()
	hkey := "maoyan_movie"
	jsonStrings, err := redisClient.HMGet(hkey, ids...).Result()

	movies := []Movie{}
	for _, item := range jsonStrings {
		if string, ok := item.(string); ok {
			movie := Movie{}
			json.Unmarshal([]byte(string), &movie)
			movies = append(movies, movie)
		}
	}

	return movies, err
}

func ParseMarkdown() error {
	tmpl, err := template.ParseFiles("template/movies") //解析模板文件

	mdFile := fmt.Sprintf("archives/movie_%s.md", time.Now().Format("2006-01-02"))

	file, err := os.Create(mdFile)
	defer file.Close()

	movies, err := FetchRedisData()
	err = tmpl.Execute(file, movies) //执行模板的merger操作
	return err
}

func FetchRedisDataHackNews() ([]HacknewsItem, error) {
	skey := time.Now().Format("hacknews-2006-01-02")
	urls, err := redisClient.SMembers(skey).Result()
	fmt.Println(urls)
	hkey := "hacknews"
	jsonStrings,err := redisClient.HMGet(hkey, urls ...).Result()

	newsItems := []HacknewsItem{}
	for idx, item := range jsonStrings {
		fmt.Println(idx)
		fmt.Println(item)

		if string, ok := item.([]byte); ok {
			items := HacknewsItem{}
			json.Unmarshal([]byte(string), &items)
			newsItems = append(newsItems, items)
		}
	}

	return newsItems, err
}
func ParseMarkdownHacknews()error{
	tmpl, err := template.ParseFiles("template/hacknews") //解析模板文件

	mdFile := fmt.Sprintf("archives/hacknews_%s.md", time.Now().Format("2006-01-02"))

	file, err := os.Create(mdFile)
	defer file.Close()

	movies, err := FetchRedisData()
	err = tmpl.Execute(file, movies) //执行模板的merger操作
	return err
}