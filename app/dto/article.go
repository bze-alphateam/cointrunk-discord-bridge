package dto

type Article struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Picture   string `json:"picture"`
	Publisher string `json:"publisher"`
	Paid      bool   `json:"paid"`
	CreatedAt string `json:"created_at"`
}

type ArticleResponse struct {
	Articles []Article `json:"article"`
}
