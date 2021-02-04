package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Thumbnail struct {
	URL string `json:"url"`
}

type Embed struct {
	Title       string    `json:"title"`
	Color       int       `json:"color"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Description string    `json:"description"`
	Footer      Footer    `json:"footer"`
	Fields      []Field   `json:"fields"`
}

type Webhook struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

func (e *Embed) SetTitle(title string) {
	e.Title = title
}

func (e *Embed) SetColor(color int) {
	e.Color = color
}

func (e *Embed) SetThumbnail(u string) {
	e.Thumbnail = Thumbnail{URL: u}
}

func (e *Embed) SetDescription(description string) {
	e.Description = description
}

func (e *Embed) SetFooter(text, icon string) {
	e.Footer = Footer{Text: text, IconURL: icon}
}

func (e *Embed) AddField(name, value string, inline bool) {
	e.Fields = append(e.Fields, Field{Name: name, Value: value, Inline: inline})
}

func (w *Webhook) SetContent(content string) {
	w.Content = content
}

func (w *Webhook) AddEmbed(e Embed) {
	w.Embeds = append(w.Embeds, e)
}

func (w *Webhook) Encode() ([]byte, error) {
	encoded, err := json.Marshal(w)

	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func (w *Webhook) Send(u string) error {

	payload, err := w.Encode()

	if err != nil {
		return err
	}

	for {
		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		resp.Body.Close()

		if resp.StatusCode != 429 {
			return nil
		}

		time.Sleep(time.Second * 5)
	}
}
