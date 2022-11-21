package main

import (
	"bot/parsing"
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Form struct {
	tg_name string 
	expediton string
	name string 
	email string
	age string
}


const adminChatID int64 = -1001602774786

func main() {
	parsing.ParseHTML("https://russiantravelgeek.com/expeditions/")
	bot, err := tgbotapi.NewBotAPI("5741027893:AAHgH5pyL7gQWm8MLyTuuG7lO9ftUvAliyQ")//защитить
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	//log.Printf("Authorized on account %s", bot.Self.UserName)
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	

	var form Form
	//type form map[string]string
	//var bd map[int64]form	
	updates := bot.GetUpdatesChan(u) 
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
            continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        if !update.Message.IsCommand() { // ignore any non-command Messages
			msg.Text = "Чтобы посмотреть список экспедиций, нажмите /find \nЧтобы подать заявку на экспедицию, нажмите /go"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			continue
        }
		
		switch update.Message.Command() {
		case "find":
			find(msg, bot)
        case "go":
			form.tg_name = update.Message.From.UserName
			ask(msg, bot, updates, &form, "Куда идем?")
			ask(msg, bot, updates, &form, "Имя и фамилия:")
			ask(msg, bot, updates, &form, "Электронный адрес:")
			ask(msg, bot, updates, &form, "Возраст:")
			//add_in_bd((bot,form))
			check_info(msg, bot, form) //todo: add ability to make changes in form/refill form
			//send_info(form)
			//https://t.me/+aCYdv4e0hmw0NWUy
        default:
            msg.Text = "I don't know that command"
        }
        }	
	}
	// func send_info(form Form) {
	// 	//https://api.telegram.org/bot5741027893:AAHgH5pyL7gQWm8MLyTuuG7lO9ftUvAliyQ/sendMessage?chat_id=[MY_CHANNEL_NAME]&text="hi"
	// 	//id -1001602774786
	// }
	func check_info(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI, form Form) {
		msg.Text = fmt.Sprintf("Ник в телеграмме : %s\nЭкспедиция : %s\nИмя : %s\nEmail : %s\nВозраст : %s", form.tg_name, form.expediton, form.name, form.email, form.age)

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		//Все правильно? Да, отправить/ Нет, заполнить заново / Не отправлять, заполню позже
		msg.ChatID = adminChatID
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
	
	func find (msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
		msg.Text = "https://russiantravelgeek.com/expeditions/"
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
	
	func ask(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, form *Form, q string) {
		msg.Text = q
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		
		for update := range updates {
			if update.Message != nil {
				//это можно сделать без повторения списков вопросов? можно. Как?
				switch q {
				case "Куда идем?": form.expediton = update.Message.Text
				case "Имя и фамилия:": form.name = update.Message.Text
				case "Электронный адрес:": form.email = update.Message.Text
				case "Возраст:": form.age = update.Message.Text
				}

				break
			}
		}
		
	}


	//todo: add go rutines for multiple use