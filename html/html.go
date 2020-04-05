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

func GenerateHTML(notes []manager.Note) (string, error) {
	html := "<html>"

	for _, note := range notes {
        html += fmt.Sprintf("<h1>%s</h1>", note.Title)
		md, err := ioutil.ReadFile(note.Path)
		if err != nil {
			return html, err
		}
		h := string(html2md.Run(md, html2md.WithNoExtensions()))
		html += incrimentHeaders(h)
	}

	html += "</html>"
	return html, nil
}
