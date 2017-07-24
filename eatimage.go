package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"fmt"

	esa "github.com/hiroakis/esa-go"
	"github.com/kaepa3/eatimage/conf"
	"github.com/kaepa3/eatimage/extraction"
)

var config conf.Config
var ConfigPath = "./config.toml"

func main() {

	config.ReadConfig(ConfigPath)

	// Initializing client
	c := esa.NewEsaClient(config.APIKey, config.TeamName)

	//get all posts
	var list []string
	page := 0
	count := 0
	for {
		c.SetPage(page)
		posts, err := c.GetPosts()
		if err != nil || len(posts.Posts) == 0 {
			break
		}
		for _, val := range posts.Posts {
			count++
			for _, val := range extraction.LinkImage(strings.NewReader(val.BodyHtml)) {
				list = append(list, val)
			}
		}
		page++
	}
	fmt.Println("count:", len(list))
	saveImage(list)
}

func saveImage(list []string) {
	for _, val := range list {
		save(val)
	}
}

func save(url string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	name := urlToFileName(url)
	fmt.Println("download:", name)
	file, err := os.Create(config.ExportRoot + `/` + name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, response.Body)
}

func urlToFileName(text string) string {
	pos := strings.LastIndex(text, "/") + 1
	return text[pos:]
}
