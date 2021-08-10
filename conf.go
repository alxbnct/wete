package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type (
	Telegram struct {
		Token string
		Socks string
		Id    int64
	}
	Conf struct {
		Telegram Telegram
	}
)

var conf *Conf

func init() {
	if _, err := toml.DecodeFile("config/config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}
