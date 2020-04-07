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
	noTagTag    = "untagged"
	notesDirKey = "$NOTES"
)

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
	html += "<input id=\"dark-mode\" type=\"checkbox\"/>"
	html += "<div class=\"tag_selector\">"
	for _, tag := range tags {
		html += fmt.Sprintf(
			"<label for=\"__id_%s\">%s</label>",
			getClass(tag),
			tag,
		)
	}
	html += "<label id=\"dark-mode-toggle\" for=\"dark-mode\">Dark mode</label>"
	html += "</div>"
	return html
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

func GenerateHTML(notes []manager.Note, notesDir string) (string, error) {
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

		bmd, err := ioutil.ReadFile(note.Path)
		if err != nil {
			return html, err
		}
		md := strings.Replace(string(bmd), notesDirKey, notesDir, -1)
		noteHtml := string(html2md.Run(
			[]byte(removeTags(md)),
			html2md.WithNoExtensions(),
		))

		html += fmt.Sprintf("<div class=\"note %s\" id=\"%s\">", classes, note.Title)
		html += "<div class=\"header\">"
		html += fmt.Sprintf("<p class=\"note-header\">%s</p>", note.Title)
		html += tagHtml
		html += "</div>"
		html += noteHtml
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

	html = getToggles(tags) + html
	return "<html>" + getStyle(tags) + "<body>" + html + "<body></html>", nil
}

func getStyle(tags []string) string {
	// We just make the CSS a big string so we can easily construct a single
	// html file that displays the notes, without relying on reading from an
	// external css file
	css := `
html {
    background-color: #F4EFE5;
    font-family: Arial, Helvetica, sans-serif;
    padding: 10px 0px;
    color: #2E2E2E;
}

body {
	margin: 0px auto;
	max-width: 800px;
}

p {
    margin: 5px 0px;
	font-size: 0.9em;
}

h1 {
    margin: 5px 0px;
    font-size: 1.6em;
}

h2 {
    margin: 5px 0px;
    font-size: 1.45em;
}

h3 {
    margin: 5px 0px;
    font-size: 1.3em;
}

h4 {
    margin: 5px 0px;
    font-size: 1.2em;
}

h5 {
    margin: 5px 0px;
    font-size: 1.1em;
}

h6 {
    margin: 5px 0px;
    font-size: 1em;
}

input {
    display: none;
}

label {
    color: #F4EFE5;
    margin: 0px 10px 0px 0px;
    padding: 3px 7px;
    border-radius: 3px;
    background-color: #6D9D99;
    white-space: nowrap;
}

label:hover {
    cursor: pointer;
}

div.tag_selector {
    overflow-x: auto;
    padding: 5px 0px;
}

#dark-mode-toggle {
    color: #2E2E2E;
    background-color: #FAF8F3;
    float: right;
    margin-right: 0px;
}

#dark-mode:checked ~ div > #dark-mode-toggle {
    color: #FAF8F3;
    background-color: #2E2E2E;
}

div.note {
    margin: 10px 0px;
    padding: 10px;
    border-radius: 3px;
    box-shadow: 0px 0px 5px grey;
    background-color: #FAF8F3;
}

div.note * {
	max-width: 100%;
}

div.note img {
	max-height: 450px;
	margin: auto;
	display: block;
	padding: 5px;
}

.note-header {
    font-weight: bold;
    margin: 1px 0px;
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
    margin: 0px 0px 0px 10px;
    padding: 1px 5px;
    border-radius: 3px;
	border: 2px solid #6D9D99;
}

div.header {
    overflow: auto;
    padding: 0px 5px 7px 0px;
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
