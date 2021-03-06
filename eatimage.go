package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"fmt"

	"github.com/BurntSushi/toml"
	esa "github.com/hiroakis/esa-go"
	"github.com/kaepa3/eatimage/conf"
	"github.com/kaepa3/eatimage/extraction"
)

type OutputFormat struct {
	Category string   `toml:"category"`
	Body     string   `toml:"body"`
	Images   []string `toml:"images"`
}
type ArayFormat struct {
	Records []OutputFormat `toml:"records"`
}

var config conf.Config
var ConfigPath = "./config.toml"

func main() {
	config.ReadConfig(ConfigPath)

	// Initializing client
	c := esa.NewEsaClient(config.APIKey, config.TeamName)

	//get all posts
	var list []string
	var bodys ArayFormat
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
			var imageList []string
			for _, text := range extraction.LinkImage(strings.NewReader(val.BodyHtml)) {
				list = append(list, text)
				imageList = append(imageList, urlToFileName(text))
			}
			info := OutputFormat{Category: val.FullName, Body: replaceChangeLine(val.BodyMd), Images: imageList}
			bodys.Records = append(bodys.Records, info)
		}
		page++
	}
	fmt.Println("count:", len(list))
	saveImage(list)
	outputToml(bodys)
}

func replaceChangeLine(text string) string {
	return text
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

func outputToml(list ArayFormat) {
	fp, err := os.OpenFile("result.toml", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()
	encoder := toml.NewEncoder(fp)
	err = encoder.Encode(list)

}
