package dto

type WebhookEmbedAuthor struct {
	Name string `json:"name"`
	Url  string `json:"url,omitempty"`
}

type WebhookEmbedImage struct {
	URL string `json:"url"`
}

type WebhookEmbedFooter struct {
	Text string `json:"text"`
}

type WebhookEmbed struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Url         string             `json:"url,omitempty"`
	Color       int                `json:"color,omitempty"`
	Image       WebhookEmbedImage  `json:"image,omitempty"`
	Author      WebhookEmbedAuthor `json:"author,omitempty"`
	Timestamp   string             `json:"timestamp,omitempty"`
	Footer      WebhookEmbedFooter `json:"footer,omitempty"`
}

type WebhookMessage struct {
	Content string         `json:"content"`
	Embeds  []WebhookEmbed `json:"embeds"`
}
