package discord

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"time"
)

// WebhookPayload is the webhook endpoint payload
type WebhookPayload struct {
	Embeds   []embedItem `json:"embeds"`
	Username string      `json:"username,omitempty"`
}

// JSON is the actual notification payload
type embedItem struct {
	Title     string       `json:"title,omitempty"`
	Content   string       `json:"description,omitempty"`
	URL       string       `json:"url,omitempty"`
	Timestamp string       `json:"timestamp,omitempty"`
	Color     uint         `json:"color,omitempty"`
	Footer    *embedFooter `json:"footer,omitempty"`
}

type embedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// CreatePayloadFromItems creates a JSON payload to be sent to the discord webhook API
func CreatePayloadFromItems(items []types.MessageItem, title string, colors [types.MessageLevelCount]uint, omitted int) (WebhookPayload, error) {

	metaCount := 1
	if omitted < 1 && len(title) < 1 {
		metaCount = 0
	}
	itemCount := util.Min(9, len(items))

	embeds := make([]embedItem, metaCount, itemCount+metaCount)

	for _, item := range items {

		color := uint(0)
		if item.Level >= types.Unknown && int(item.Level) < len(colors) {
			color = colors[item.Level]
		}

		ei := embedItem{
			Content: item.Text,
			Color:   color,
		}

		if item.Level != types.Unknown {
			ei.Footer = &embedFooter{
				Text: item.Level.String(),
			}
		}

		if !item.Timestamp.IsZero() {
			ei.Timestamp = item.Timestamp.UTC().Format(time.RFC3339)
		}

		embeds = append(embeds, ei)
	}

	embeds[0].Title = title
	if omitted > 0 {
		embeds[0].Footer = &embedFooter{
			Text: fmt.Sprintf("... (%v character(s) where omitted)", omitted),
		}
	}

	return WebhookPayload{
		Embeds: embeds,
	}, nil
}
