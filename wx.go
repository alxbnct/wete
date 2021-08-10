package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/skip2/go-qrcode"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func wxInit() (*openwechat.Bot, *openwechat.Self, error) {
	bot := openwechat.DefaultBot()
	bot.UUIDCallback = ConsoleQrCode

	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	err := bot.HotLogin(reloadStorage, true)
	if err != nil {
		fmt.Println(err)
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	checkErr(err)

	return bot, self, err

}

func wxGet() {

	bot, self, err := wxInit()
	checkErr(err)
	friends, err := self.Friends()
	checkErr(err)
	/* 	groups, err := self.Groups()
	   	checkErr(err) */

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.IsSendByFriend() && len(msg.Content) != 0 {

			//fmt.Println(msg.Content)

			sender, err := msg.Sender()
			if err != nil {
				fmt.Println(err)
			}
			if sender.RemarkName == "neo" {
				fmt.Println(msg.Content)

			}

			bot, updates, err := tgInit()
			if err != nil {
				fmt.Println(err)
			}

			xx := tg.NewMessage(conf.Telegram.Id, msg.Content)
			bot.Send(xx)

			for update := range updates {
				// ignore any non-Message Updates
				if update.Message == nil && update.Message.Chat.ID != conf.Telegram.Id {
					continue
				}

				ok := friends.SearchByRemarkName(1, "neo")

				if ok.Count() > 0 {
					go wxSend(ok.First(), update.Message.Text)
				}

			}
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func wxSend(f *openwechat.Friend, msg string) {
	f.SendText(msg)
}

type Mee struct {
	ChatID             int64
	SuperGroupUsername string
}
