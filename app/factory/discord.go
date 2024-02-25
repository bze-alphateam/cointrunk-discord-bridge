package factory

import (
	"fmt"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/dto"
	"strconv"
	"time"
)

const (
	maxDiscordTitleLen       = 256
	longTitleSuffix          = "..."
	paidArticleColorCode     = 15844367
	paidArticleDefaultAuthor = "PAID ARTICLE"

	publisherUrl = "https://app.cointrunk.io/publisher/%s"

	footerText            = "You are viewing permissionless content submitted on BeeZee (BZE) blockchain. All the content is submitted under blockchain's own rules. Visit CoinTrunk.io for more details."
	footerPaidArticleText = " THIS ARTICLE IS SPONSORED!"
)

type PublisherStorage interface {
	GetPublisherDetails(publisherAddress string) *dto.Publisher
}

type WebhookFactory struct {
	publisherStorage PublisherStorage
}

func NewWebhookFactory(publisherStorage PublisherStorage) (*WebhookFactory, error) {
	if publisherStorage == nil {
		return nil, fmt.Errorf("invalid dependencies provided to factory constructor")
	}

	return &WebhookFactory{publisherStorage: publisherStorage}, nil
}

func (wf WebhookFactory) NewWebhookMessage(article dto.Article) (*dto.WebhookMessage, error) {
	createdAtInt, err := strconv.ParseInt(article.CreatedAt, 10, 64)
	if err != nil {
		return nil, err
	}

	t := time.Unix(createdAtInt, 0)
	//discord requires certain format
	discordTime := t.Format(time.RFC3339)
	title := getTrimmedTitle(article.Title)

	//basic embed article
	articleEmbed := dto.WebhookEmbed{
		Title:     title,
		Url:       article.Url,
		Timestamp: discordTime,
	}

	//check if URL exists and add it
	if len(article.Picture) > 0 {
		articleEmbed.Image = dto.WebhookEmbedImage{URL: article.Picture}
	}

	author := dto.WebhookEmbedAuthor{
		Name: paidArticleDefaultAuthor,
	}

	publisherDetails := wf.publisherStorage.GetPublisherDetails(article.Publisher)
	if publisherDetails != nil {
		author.Url = fmt.Sprintf(publisherUrl, publisherDetails.Address)
		author.Name = publisherDetails.Name
	}
	articleEmbed.Author = author

	fText := footerText
	if article.Paid {
		articleEmbed.Color = paidArticleColorCode
		fText += footerPaidArticleText
	}

	articleEmbed.Footer = dto.WebhookEmbedFooter{Text: fText}

	return &dto.WebhookMessage{
		Embeds: []dto.WebhookEmbed{articleEmbed},
	}, nil
}

// getTrimmedTitle - discord accepts a title of max length maxDiscordTitleLen
// This function returns the title trimmed if it exceeds the limit discord imposes + a suffix at the end, ex: ...
func getTrimmedTitle(title string) string {
	if len(title) > maxDiscordTitleLen {
		title = title[0:(maxDiscordTitleLen - 1 - len(longTitleSuffix))]
		title += longTitleSuffix
	}

	return title
}
