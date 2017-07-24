package extraction

import (
	"fmt"
	"io"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

func LinkImage(sr io.Reader) []string {
	doc, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	list := make([]string, 0, 2000)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if isImage(url) {
			list = append(list, url)
		}
	})
	return list
}

func isImage(text string) bool {
	list := []string{".png", ".jpeg", ".jpg"}
	for _, val := range list {
		if strings.Contains(text, val) {
			return true
		}
	}
	return false
}
