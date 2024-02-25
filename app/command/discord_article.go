package command

import (
	"fmt"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/dto"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/service"
	"log"
	"strconv"
)

type ArticleListener interface {
	ListenForArticles(handler service.MessageHandler) error
}

type WebhookMessageCreator interface {
	NewWebhookMessage(article dto.Article) (*dto.WebhookMessage, error)
}

type ArticleProvider interface {
	FetchArticles() ([]dto.Article, error)
}

type LastIdProvider interface {
	GetLastInsertedID() int64
	SaveLastInsertedID(id int64) error
}

type DiscordClient interface {
	PostMessage(message dto.WebhookMessage) error
}

type DiscordArticle struct {
	articleProvider ArticleProvider
	idProvider      LastIdProvider
	messageCreator  WebhookMessageCreator
	client          DiscordClient
	listener        ArticleListener
}

func NewDiscordArticle(articleProvider ArticleProvider, idProvider LastIdProvider, messageCreator WebhookMessageCreator, client DiscordClient, listener ArticleListener) (*DiscordArticle, error) {
	if articleProvider == nil || idProvider == nil || messageCreator == nil || client == nil || listener == nil {
		return nil, fmt.Errorf("invalid dependencies provider to DiscordArticle command")
	}

	return &DiscordArticle{
		articleProvider: articleProvider,
		idProvider:      idProvider,
		messageCreator:  messageCreator,
		client:          client,
		listener:        listener,
	}, nil
}

func (da *DiscordArticle) Start() error {
	da.postArticles()
	return da.listener.ListenForArticles(da.messageHandler())
}

func (da *DiscordArticle) messageHandler() service.MessageHandler {
	return func(msg []byte) {
		da.postArticles()
	}
}

func (da *DiscordArticle) postArticles() {
	articles, err := da.articleProvider.FetchArticles()
	if err != nil {
		log.Println(fmt.Sprintf("could not fetch articles: %v", err))
		return
	}

	lastId := da.idProvider.GetLastInsertedID()
	articles = da.filterArticlesById(articles, lastId)
	log.Println(fmt.Sprintf("articles found to be published: %d", len(articles)))

	da.sendToDiscord(articles)
}

func (da *DiscordArticle) sendToDiscord(articles []dto.Article) {
	for i := len(articles) - 1; i >= 0; i-- {
		msg, err := da.messageCreator.NewWebhookMessage(articles[i])
		if err != nil {
			log.Println(fmt.Sprintf("error when creating webhook message: %v", err))
			continue
		}

		err = da.client.PostMessage(*msg)
		if err != nil {
			log.Println(fmt.Sprintf("error when sending webhook: %v", err))
			continue
		}

		artId := da.convertArticleId(articles[i].Id)
		if artId == 0 {
			continue
		}

		err = da.idProvider.SaveLastInsertedID(artId)
		if err != nil {
			log.Println(fmt.Sprintf("error when saving last id: %v", err))
			continue
		}
	}
}

func (da *DiscordArticle) filterArticlesById(articles []dto.Article, lastId int64) []dto.Article {
	var filtered []dto.Article
	for _, art := range articles {
		artId := da.convertArticleId(art.Id)
		if artId == 0 {
			continue
		}

		if artId <= lastId {
			break
		}

		filtered = append(filtered, art)
		continue
	}

	return filtered
}

func (da *DiscordArticle) convertArticleId(id string) int64 {
	artId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(fmt.Sprintf("error when converting article id to int64: %v", err))
		return 0
	}

	return artId
}
