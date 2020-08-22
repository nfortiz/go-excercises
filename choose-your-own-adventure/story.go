package cyoa

import (
	"encoding/json"
	"io"
)

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
	Arc  string `json:"arc"`
}

type Demo struct {}