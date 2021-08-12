package main

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/eatmoreapple/openwechat"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/skip2/go-qrcode"
)

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func start() {
	bot := openwechat.DefaultBot()

	var count int32
	bot.GetMessageErrorHandler = func(err error) {
		atomic.AddInt32(&count, 1)
		if count == 3 {
			bot.Logout()
		}
	}

	bot.UUIDCallback = ConsoleQrCode

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	bot.HotLogin(reloadStorage, true)
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("当前登录用户：", user)
	}

	// 获取所有的好友
	friends, err := user.Friends()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("获取通讯录中好友列表成功共", friends.Count(), "个")
	}

	// 获取所有的群组
	groups, err := user.Groups()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("获取通讯录中群聊列表成功，共", groups.Count(), "个")
	}

	mps, err := user.Mps()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("获取通讯录中公众号列表成功，共", mps.Count(), "个")
	}

	fh, err := user.FileHelper()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("获取文件传输助手成功", fh.UserName, fh.RemarkName)
	}

	bott, updates, err := tgInit()
	if err != nil {
		fmt.Println(err)
	}

	for update := range updates {
		// ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		bot.MessageHandler = func(msg *openwechat.Message) {
			if len(msg.Content) != 0 && msg.IsSendByGroup() {
				a, b, c := TextMessageHandler(msg)
				if a == "国光帮帮忙" {
					log.Printf("[%v]%v: %v", a, b, c)
					cc := fmt.Sprintf("[%v]%v: %v", a, b, c)
					xx := tg.NewMessage(update.Message.Chat.ID, cc)
					bott.Send(xx)
				}

			}
		}
	}

	bot.Block()
}

func TextMessageHandler(c *openwechat.Message) (string, string, string) {
	sender, _ := c.Sender()
	senderUser := sender.NickName
	SenderInGroup, _ := c.SenderInGroup()
	return senderUser, SenderInGroup.NickName, c.Content
}
