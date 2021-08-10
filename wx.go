package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

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
	if err != nil {
		fmt.Println(err)
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	//bot.Block()

	return bot, self, err

}

func wxGet() {

	bot, self, err := wxInit()
	if err != nil {
		fmt.Println(err)
	}
	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && len(msg.Content) != 0 {
			/* msg.ReplyText("pong") */
			/* fmt.Println(msg.Content) */
			sender, err := msg.Sender()
			if err != nil {
				fmt.Println(err)
			}
			if sender.RemarkName == "neo" {
				fmt.Println(msg.Content)
			}
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
