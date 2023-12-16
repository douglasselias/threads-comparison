package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func renderText(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text += n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += renderText(c)
	}
	return text
}

func getTextByDataPid(url, dataPid string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return "", err
	}

	var text string
	var search func(*html.Node)
	search = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == dataPid {
					text += strings.TrimSpace(html.UnescapeString(renderText(n)))
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}

	search(doc)
	return text, nil
}

func main() {
	// response, err := http.Get("https://www.jw.org/pt/biblioteca/jw-apostila-do-mes/janeiro-fevereiro-2024-mwb/Programa%C3%A7%C3%A3o-da-Reuni%C3%A3o-Vida-e-Minist%C3%A9rio-para-1-%E2%81%A07-de-janeiro-de-2024/")
	// if err != nil {
	// 	fmt.Println("Error requesting webpage: ", err)
	// 	return
	// }
	// defer response.Body.Close()

	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body: ", err)
	// 	return
	// }

	text, err := getTextByDataPid("https://www.jw.org/pt/biblioteca/jw-apostila-do-mes/janeiro-fevereiro-2024-mwb/Programa%C3%A7%C3%A3o-da-Reuni%C3%A3o-Vida-e-Minist%C3%A9rio-para-1-%E2%81%A07-de-janeiro-de-2024/", "p5")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Text inside div with data-pid:", text)
	// fmt.Println("Body: ", string(body))
}
