package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/BurntSushi/toml"
	"github.com/eatmoreapple/openwechat"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/skip2/go-qrcode"
	"golang.org/x/net/proxy"
)

type (
	Telegram struct {
		Token string
		Socks string
		Id    int64
		Debug bool
	}

	QQ struct {
	}
	Conf struct {
		Telegram Telegram
		QQ       QQ
	}
)

var conf *Conf

func init() {
	if _, err := toml.DecodeFile("config/config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

func ProxyNewBotApi(token string) (*tg.BotAPI, error) {
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

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func GroupMessageHandler(c *openwechat.Message) (string, string, string) {
	sender, _ := c.Sender()
	senderUser := sender.NickName
	SenderInGroup, _ := c.SenderInGroup()
	return senderUser, SenderInGroup.NickName, c.Content
}

func FriendMessageHandler(c *openwechat.Message) (string, string) {
	sender, _ := c.Sender()
	return sender.RemarkName, c.Content
}

func main() {
	wx := openwechat.DefaultBot()

	var count int32
	wx.GetMessageErrorHandler = func(err error) {
		atomic.AddInt32(&count, 1)
		if count == 3 {
			wx.Logout()
		}
	}

	wx.UUIDCallback = ConsoleQrCode

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	wx.HotLogin(reloadStorage, true)
	user, err := wx.GetCurrentUser()
	if err != nil {
		log.Println(err)
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

	bot, err := ProxyNewBotApi(conf.Telegram.Token)
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

	for update := range updates {
		if update.Message == nil {
			continue
		}

		wx.MessageHandler = func(msg *openwechat.Message) {
			if len(msg.Content) != 0 && msg.IsText() {
				a, b, c := GroupMessageHandler(msg)
				d, e := FriendMessageHandler(msg)
				if a == "我不是鸽手" {
					log.Printf("[%v]%v: %v", a, b, c)
					cc := fmt.Sprintf("[%v]%v: %v", a, b, c)
					xx := tg.NewMessage(update.Message.Chat.ID, cc)
					bot.Send(xx)
				}

				if a == "国光帮帮忙" {
					log.Printf("[%v]%v: %v", a, b, c)
					cc := fmt.Sprintf("[%v]%v: %v", a, b, c)
					xx := tg.NewMessage(update.Message.Chat.ID, cc)
					bot.Send(xx)
				}

				if d == "xx" {
					log.Printf("%v: %v", d, e)
					cc := fmt.Sprintf("%v: %v", d, e)
					xx := tg.NewMessage(update.Message.Chat.ID, cc)
					bot.Send(xx)
				}

				if d == "neo" {
					log.Printf("%v: %v", d, e)
					cc := fmt.Sprintf("%v: %v", d, e)
					xx := tg.NewMessage(update.Message.Chat.ID, cc)
					bot.Send(xx)
				}

			}
		}
	}
	wx.Block()
}
