package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bze-alphateam/cointrunk-discord-bridge/app/dto"
)

const (
	articleUrl   = "/bze/cointrunk/all_articles?pagination.count_total=false&pagination.limit=50&pagination.reverse=true"
	publisherUrl = "/bze/cointrunk/publisher/%s"
)

type RestDataProvider struct {
	restUrl string

	publishers map[string]*dto.Publisher
}

func NewRestDataProvider(restUrl string) (*RestDataProvider, error) {
	if restUrl == "" {
		return nil, fmt.Errorf("invalid rest url provided to constructor")
	}

	return &RestDataProvider{
		restUrl:    restUrl,
		publishers: make(map[string]*dto.Publisher),
	}, nil
}

func (p *RestDataProvider) FetchArticles() ([]dto.Article, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", p.restUrl, articleUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var articles dto.ArticleResponse
	err = json.Unmarshal(body, &articles)
	if err != nil {
		return nil, err
	}

	return articles.Articles, nil
}

func (p *RestDataProvider) GetPublisherDetails(publisherAddress string) *dto.Publisher {
	publisher, found := p.publishers[publisherAddress]
	if found && publisher != nil {
		return publisher
	}

	publisher, err := p.fetchPublisher(publisherAddress)
	if err != nil {
		return nil
	}

	p.publishers[publisherAddress] = publisher

	return publisher
}

func (p *RestDataProvider) fetchPublisher(publisherAddress string) (*dto.Publisher, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", p.restUrl, p.getPublisherUrl(publisherAddress)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pResp dto.PublisherResponse
	err = json.Unmarshal(body, &pResp)
	if err != nil {
		return nil, err
	}

	return pResp.Publisher, nil
}

func (p *RestDataProvider) getPublisherUrl(publisherAddress string) string {
	return fmt.Sprintf(publisherUrl, publisherAddress)
}
