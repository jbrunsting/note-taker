package html

import (
	"fmt"
	"io/ioutil"
	"sort"
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

func getClass(tag string) string {
	o := ""
	for _, c := range tag {
		if ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') {
			o += string(c)
		} else if c == ' ' {
			o += "_"
		}
	}
	return o
}

func getToggles(tags []string) string {
	html := ""
	for _, tag := range tags {
		html += fmt.Sprintf(
			"<input id=\"__id_%s\" class=\"%s\" type=\"checkbox\"/>",
			getClass(tag),
			getClass(tag),
		)
		html += fmt.Sprintf(
			"<label for=\"__id_%s\">%s</label>",
			getClass(tag),
			tag,
		)
	}
	return html
}

func getStyle(tags []string) string {
	css := `
input {
    display: none;
}

label {
    color: white;
    margin: 5px;
    padding: 3px 7px;
    border-radius: 3px;
    background-color: #0080ff;
}

input:checked+label {
    background-color: #73b9ff;
}

label:hover {
    cursor: pointer;
}
`
	for _, tag := range tags {
		css += fmt.Sprintf(
			"input.%[1]s:checked ~ .%[1]s {display:none}",
			getClass(tag),
		)
	}
	return "<style>" + css + "</style>"
}

type OrderedTag struct {
	Tag   string
	Count int
}

func GenerateHTML(notes []manager.Note) (string, error) {
	oTags := make(map[string]*OrderedTag)
	html := ""
	for _, note := range notes {
		classes := ""
		for _, tag := range note.Tags {
			tag = strings.ToLower(tag)

			if _, ok := oTags[tag]; !ok {
				oTags[tag] = &OrderedTag{tag, 0}
			}

			oTags[tag].Count += 1
			classes += " " + getClass(tag)
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

	vals := []*OrderedTag{}
	for _, v := range oTags {
		vals = append(vals, v)
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].Count > vals[j].Count
	})

	tags := []string{}
	for _, ot := range vals {
		tags = append(tags, ot.Tag)
	}

	html = getStyle(tags) + getToggles(tags) + html

	return "<html>" + html + "</html>", nil
}
