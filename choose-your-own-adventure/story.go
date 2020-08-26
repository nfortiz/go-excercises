package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var defaultTmpl string = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Choose Your Own Adventure</title>
	<style>
		body { font-family: helvetica, arial; }
		h1 {
			text-align: center;	
			position: relative;	
		}
		.page {
			width: 80%;
			max-width: 500px;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
			padding: 80px;
			background: #FFFCF6;
			border: 1px solid #eee;
			box-shadow: 0 10px 6px -6px #777;
		}
		ul {
			border-top: 1px dotted #ccc;
			padding: 10px 0 0 0;
		}
		li {
			padding-top: 10px;
		}
	</style>
</head>
<body>
	<section class="page">
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
	</section>
</body>
</html>
`
var tpl = template.Must(template.New("").Parse(defaultTmpl))

type HandlerOption func (h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func (h *handler) {
		h.t = t
	}
}

func NewHandler (s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {

		err := h.t.Execute(w, chapter)

		if err != nil {
			http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func JSONStory(r io.Reader) (Story, error) {
	decoded := json.NewDecoder(r)
	var story Story
	if err := decoded.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Paragraphs   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter  string `json:"arc"`
}

type Demo struct {}