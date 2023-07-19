package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
)

// TODO: Config to change values

func formatDuration(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second

	var result string
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 || (hours == 0 && minutes == 0) {
		result += fmt.Sprintf("%ds", seconds)
	}

	return result
}

func NewDiscordBlocks(url string, network string, upgradeTime time.Duration, currentHeight uint64) {
	upgradeTimeStr := formatDuration(upgradeTime)

	embed := discordwebhook.Embed{
		Title:     network + " BLOCKS",
		Color:     46628,
		Timestamp: time.Now(),
		Thumbnail: discordwebhook.Thumbnail{
			Url: "https://st2.depositphotos.com/4431055/11855/i/600/depositphotos_118559772-stock-photo-color-childish-block-to-the.jpg",
		},
		Fields: []discordwebhook.Field{
			{
				Name:   "Time taken to Upgrade",
				Value:  upgradeTimeStr,
				Inline: true,
			},
			{
				Name:   "Current Height",
				Value:  fmt.Sprintf("%d", currentHeight),
				Inline: false,
			},
		},
		Footer: discordwebhook.Footer{
			Text: "blocks by reece",
		},
	}

	SendEmbed(url, "https://miro.medium.com/v2/resize:fit:2400/1*HJxtRYhZYs7vhWH25hmWHA.png", embed)
}

func NewDiscordTimeToUpgrade(url string, network string, currentHeight uint64) {

	embed := discordwebhook.Embed{
		Title:     network + " Time to Upgrade...",
		Color:     16776960,
		Timestamp: time.Now(),
		Thumbnail: discordwebhook.Thumbnail{
			Url: "https://static-00.iconduck.com/assets.00/pending-icon-512x504-9zrlrc78.png",
		},
		Fields: []discordwebhook.Field{
			{
				Name:   "Current Height",
				Value:  fmt.Sprintf("%d", currentHeight),
				Inline: true,
			},
		},
		Footer: discordwebhook.Footer{
			Text: "pending by reece",
		},
	}

	SendEmbed(url, "https://miro.medium.com/v2/resize:fit:2400/1*HJxtRYhZYs7vhWH25hmWHA.png", embed)
}

func SendEmbed(link string, avatar string, embeds discordwebhook.Embed) error {
	hook := discordwebhook.Hook{
		Username:   "Captain Hook",
		Avatar_url: avatar,
		// Content:    "Message",
		Embeds: []discordwebhook.Embed{embeds},
	}

	payload, err := json.Marshal(hook)
	if err != nil {
		log.Fatal(err)
	}
	err = discordwebhook.ExecuteWebhook(link, payload)
	return err
}
