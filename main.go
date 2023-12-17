package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
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

func parseBody(body io.Reader) ([]map[string]string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	var currentResult map[string]string

	var search func(*html.Node)
	search = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "h3" || n.Data == "p" || n.Data == "h2" || n.Data == "h1") {
			text := strings.TrimSpace(html.UnescapeString(renderText(n)))

			if n.Data == "h1" {
				for _, attr := range n.Attr {
					if attr.Key == "id" && attr.Val == "p1" {
						currentResult = make(map[string]string)
						results = append(results, currentResult)
						currentResult[n.Data] = text
					}
				}
			} else if n.Data == "h2" {
				for _, attr := range n.Attr {
					if attr.Key == "id" && attr.Val == "p2" {
						currentResult = make(map[string]string)
						results = append(results, currentResult)
						currentResult[n.Data] = text
					}
				}
			} else if n.Data == "h3" {
				// If current node is h3, create a new result entry
				currentResult = make(map[string]string)
				results = append(results, currentResult)
				currentResult[n.Data] = text
			} else if currentResult != nil {
				// If currentResult is not nil, update the paragraph entry
				currentResult["p"] = text
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}

	search(doc)
	return results, nil
}

func main() {
	// @todo: generate url on request
	// year := 2024
	// monthTuple := (jan, fev)
	// month := jan
	// week := (1 7)
	url := "https://www.jw.org/pt/biblioteca/jw-apostila-do-mes/janeiro-fevereiro-2024-mwb/Programa%C3%A7%C3%A3o-da-Reuni%C3%A3o-Vida-e-Minist%C3%A9rio-para-1-%E2%81%A07-de-janeiro-de-2024/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	results, err := parseBody(response.Body)
	if err != nil {
		fmt.Println("Error on parsing body:", err)
		return
	}

	type Meeting struct {
		Week          string
		BibleText     string
		StartSong     string
		MainDiscourse string
		Ministry      []string
		MiddleSong    string
		ChristianLife []string
		EndSong       string
	}

	meeting := Meeting{
		Week:          "",
		BibleText:     "",
		StartSong:     "",
		MainDiscourse: "",
		Ministry:      []string{},
		MiddleSong:    "",
		ChristianLife: []string{},
		EndSong:       "",
	}

	afterMiddleSong := false

	for _, result := range results {
		if result["h1"] != "" {
			meeting.Week = result["h1"]
		}
		if result["h2"] != "" {
			meeting.BibleText = result["h2"]
		}
		// meeting.BibleText = result["h2"]

		title := result["h3"]
		time := result["p"]
		parts := strings.Split(title, ".")
		num, err := strconv.Atoi(parts[0])

		if err == nil && num > 0 && num < 10 {
			if num == 1 {
				meeting.MainDiscourse = title
			} else if num > 3 && !afterMiddleSong {
				parts := strings.Split(time, ")")
				meeting.Ministry = append(meeting.Ministry, title+" "+parts[0]+")")
			} else if num > 3 && afterMiddleSong && !strings.Contains(title, "Estudo bíblico de congregação") {
				parts := strings.Split(time, ")")
				meeting.ChristianLife = append(meeting.ChristianLife, title+" "+parts[0]+")")
			}
		}

		hasStartSong := strings.Contains(title, "Comentários iniciais")
		hasEndSong := strings.Contains(title, "Comentários finais")
		hasMiddleSong := strings.Contains(title, "Cântico")

		if err != nil && (hasStartSong || hasMiddleSong || hasEndSong) {
			if hasStartSong {
				parts := strings.Split(title, " ")
				meeting.StartSong = parts[1]
			} else if hasEndSong {
				parts := strings.Split(title, "Cântico ")
				parts1 := strings.Split(parts[1], " ")
				meeting.EndSong = parts1[0]
			} else if hasMiddleSong {
				parts := strings.Split(title, "Cântico")
				song := ""
				for _, char := range parts[1] {
					if char != 160 {
						song += string(char)
					}
				}
				meeting.MiddleSong = song
				afterMiddleSong = true
			}
		}
	}

	fmt.Println(meeting)

	/// ---------- server ------------ ///
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template, passing the Page data.
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("pdf.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template, passing the Page data.
		err = tmpl.Execute(w, meeting)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the server on port 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	http.ListenAndServe(addr, nil)
}
