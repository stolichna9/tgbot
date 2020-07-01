package main

type Config struct {
	TelegramToken string `json:"telegramToken"`
	QiwiToken     string `json:"qiwiToken"`
	QiwiWallet    string `json:"qiwiWallet"`
}
