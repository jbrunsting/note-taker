package html

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/jbrunsting/note-taker/manager"
	html2md "github.com/russross/blackfriday/v2"
)

const (
	noTagTag = "untagged"
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
	}
	html += "<div class=\"tag_selector\">"
	for _, tag := range tags {
		html += fmt.Sprintf(
			"<label for=\"__id_%s\">%s</label>",
			getClass(tag),
			tag,
		)
	}
	html += "</div>"
	return html
}

func getStyle(tags []string) string {
	// We just make the CSS a big string so we can easily construct a single
	// html file that displays the notes, without relying on reading from an
	// external css file
	css := `
html {
    background-color: #F4EFE5;
    font-family: Arial, Helvetica, sans-serif;
    padding: 10px;
    color: #2E2E2E;
}

input {
    display: none;
}

label {
    color: #F4EFE5;
    margin: 5px;
    padding: 3px 7px;
    border-radius: 3px;
    background-color: #6D9D99;
    white-space: nowrap;
}

div.tag_selector {
    overflow-x: auto;
    padding: 5px;
}

label:hover {
    cursor: pointer;
}

div.note {
    margin: 10px 10px;
    padding: 10px;
    border-radius: 3px;
    box-shadow: 0px 0px 5px grey;
    background-color: #FAF8F3;
}

h1.note-header {
    margin: 5px 0px;
    font-size: 1em;
    display: inline-block;
}

div.tag {
	display: inline-block;
	float: right;
    margin: 0px -5px;
    font-size: 0.8em;
}

div.tag > p {
    font-size: 1em;
    display: inline-block;
    margin: 5px;
    padding: 1px 5px;
    border-radius: 3px;
	border: 2px solid #6D9D99;
}

div.header {
    overflow: auto;
    padding: 0px 5px 5px 4px;
	border-bottom: 1px solid #2E2E2E;
}
`
	for _, tag := range tags {
		css += fmt.Sprintf(`
input.%[1]s ~ div.%[1]s {
    display:none
}
input.%[1]s:not(:checked) ~ div.%[1]s {
	display:block;
}
input.%[1]s:checked ~ div > label[for=__id_%[1]s] {
	background-color: #BFC9BC
}
`,
			getClass(tag),
		)
	}
	return "<style>" + css + "</style>"
}

func removeTags(md string) string {
	lines := strings.SplitN(md, "\n", 2)
	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) > 0 && firstLine[0] == '[' && firstLine[len(firstLine)-1] == ']' {
		if len(lines) > 1 {
			return lines[1]
		}
		return ""
	}
	return md
}

type OrderedTag struct {
	Tag   string
	Count int
}

func GenerateHTML(notes []manager.Note) (string, error) {
	oTags := make(map[string]*OrderedTag)
	html := ""
	for _, note := range notes {
		tagHtml := "<div class=\"tag\">"
		classes := ""
		for _, tag := range note.Tags {
			tag = strings.ToLower(tag)

			if _, ok := oTags[tag]; !ok {
				oTags[tag] = &OrderedTag{tag, 0}
			}

			oTags[tag].Count += 1
			classes += " " + getClass(tag)

			tagHtml += fmt.Sprintf("<p>%s</p>", tag)
		}
		if len(note.Tags) == 0 {
			oTags[noTagTag] = &OrderedTag{noTagTag, 0}
			classes += " " + getClass(noTagTag)
		}
		tagHtml += "</div>"

		md, err := ioutil.ReadFile(note.Path)
		if err != nil {
			return html, err
		}
		noteHtml := string(html2md.Run(
			[]byte(removeTags(string(md))),
			html2md.WithNoExtensions(),
		))

		html += fmt.Sprintf("<div class=\"note %s\">", classes)
		html += "<div class=\"header\">"
		html += fmt.Sprintf("<h1 class=\"note-header\">%s</h1>", note.Title)
		html += tagHtml
		html += "</div>"
		html += incrimentHeaders(noteHtml)
		html += "</div>"
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
