package dto

type Publisher struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Active        bool   `json:"active"`
	ArticlesCount int    `json:"articles_count"`
	CreatedAt     string `json:"created_at"`
	Respect       string `json:"respect"`
}

type PublisherResponse struct {
	Publisher *Publisher `json:"publisher"`
}
