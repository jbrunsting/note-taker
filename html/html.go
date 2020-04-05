package html

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jbrunsting/note-taker/manager"
	html2md "github.com/russross/blackfriday/v2"
)

func incrimentHeaders(html string) string {
	for i := 5; i > 0; i-- {
		o := fmt.Sprintf("<h%d>", i)
		c := fmt.Sprintf("</h%d>", i)
		no := fmt.Sprintf("<h%d>", i+1)
		nc := fmt.Sprintf("</h%d>", i+1)
		html = strings.Replace(html, o, no, -1)
		html = strings.Replace(html, c, nc, -1)
	}
	return html
}

func makeAlphanumeric(s string) string {
	o := ""
	for _, c := range s {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
			o += string(c)
		}
	}
	return o
}

func GenerateHTML(notes []manager.Note) (string, error) {
	html := "<html>"

	for _, note := range notes {
		classes := ""
		for _, tag := range note.Tags {
			classes += " " + makeAlphanumeric(tag)
		}

		html += fmt.Sprintf("<div class=\"%s\">", classes)
		html += fmt.Sprintf("<h1>%s</h1>", note.Title)
		md, err := ioutil.ReadFile(note.Path)
		if err != nil {
			return html, err
		}
		h := string(html2md.Run(md, html2md.WithNoExtensions()))
		html += incrimentHeaders(h)
		html += fmt.Sprintf("</div>")
	}

	html += "</html>"
	return html, nil
}
