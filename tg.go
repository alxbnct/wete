package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

func NewBotProxyApi(token string) (*tg.BotAPI, error) {
	client := &http.Client{}
	socks5 := "socks5://" + conf.Telegram.Socks
	if len(socks5) > 0 {
		tgProxyURL, err := url.Parse(socks5)
		if err != nil {
			log.Printf("Failed to parse proxy URL:%s\n", err)
			return nil, err
		}
		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			log.Printf("Failed to obtain proxy dialer: %s\n", err)
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}
	return tg.NewBotAPIWithClient(token, client)
}

func tgInit() (*tg.BotAPI, tg.UpdatesChannel, error) {
	bot, err := NewBotProxyApi(conf.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
	}

	return bot, updates, err
}
