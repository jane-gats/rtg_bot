package main

import (
	"bot/parsing"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Form struct {
	tg_name    string
	expediton  string
	name       string
	email      string
	age        string
	social     string
	experience string
	work       string
	speach     string
	proud      string
	companion  string
	restricts  string
	conversion string
	agreement  string
	relative   string
}

const (
	NilKeyboard              = 0
	YesNoKeyboard            = 1
	ChooseExpeditionKeyboard = 2
)

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		User struct {
			Username string `json:"username"`
		} `json:"from"`
	} `json:"message"`
}

type inlineKeyboardMarkup struct {
	AllButtons [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type inlineKeyboardButton struct {
	Text         string `json:"text,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage

type sendMessageReqBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
	ReplyMarkup *inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

const adminChatID int64 = -1001602774786
const expeditionUrl = "https://russiantravelgeek.com/expeditions/"

const apiString = "5741027893:AAHgH5pyL7gQWm8MLyTuuG7lO9ftUvAliyQ" //защитить

var openForms = make(map[int64]*Form, 100)
var expeditionList []parsing.Expedition

func main() {
	fmt.Println("Started")
	expeditionList = parsing.FetchExpeditionsFromUrl(expeditionUrl)
	go continiousFetchExpedition()

	http.ListenAndServe(":3000", http.HandlerFunc(menuHandler))

}

func menuHandler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	var text = body.Message.Text
	if strings.Contains(strings.ToLower(text), "/") {
		switch text {
		case "/find":
			if err := find(body.Message.Chat.ID); err != nil {
				fmt.Println("error in sending reply:", err)
			}
		case "/go":
			if err := ask(body.Message.Chat.ID, body.Message.User.Username); err != nil {
				fmt.Println("error in sending reply:", err)
			}
		default:
			return
		}
	} else {
		var chatId = body.Message.Chat.ID
		if _, ok := openForms[chatId]; !ok {
			sendMessageReq("Чтобы посмотреть список экспедиций, нажмите /find \nЧтобы подать заявку на экспедицию, нажмите /go", body.Message.Chat.ID, NilKeyboard)
			return
		}
		v := openForms[chatId]
		if v.tg_name == "" {
			sendMessageReq("!Чтобы посмотреть список экспедиций, нажмите /find \nЧтобы подать заявку на экспедицию, нажмите /go", body.Message.Chat.ID, NilKeyboard)
			return
		}
		switch {
		case v.expediton == "":
			v.expediton = text
			fmt.Println("text = " + text)
			sendMessageReq("Имя и Фамилия", chatId, NilKeyboard)
		case v.name == "":
			v.name = text
			sendMessageReq("Электронный адрес", chatId, NilKeyboard)
		case v.email == "":
			v.email = text
			sendMessageReq("Возраст", chatId, NilKeyboard)
		case v.age == "":
			v.age = text
			sendMessageReq("Ссылка на самую “обжитую” социальную сеть", chatId, NilKeyboard)
		case v.social == "":
			v.social = text
			sendMessageReq("Мы ранее путешествовали вместе? Есть ли опыт экспедиций с RTG", chatId, YesNoKeyboard)
			// да/нет
		case v.experience == "":
			v.experience = text
			sendMessageReq("Кем ты работаешь и чем занимаешься в обычной жизни? Можно без особого уточнения. Просто интересно узнать тебя получше", chatId, NilKeyboard)
		case v.work == "":
			v.work = text
			sendMessageReq("Есть желание провести мастер-класс или выступить с интересным спичем? Если да - назови тему", chatId, NilKeyboard)
		case v.speach == "":
			v.speach = text
			sendMessageReq("Чем ты гордишься? Один или пара моментов из жизни, вспоминая которые, ты испытываешь чувство гордости или удовлетворения", chatId, NilKeyboard)
		case v.proud == "":
			v.proud = text
			sendMessageReq("Ты едешь одна/один или с кем-то? Если да - укажи фамилию, чтобы мы рассмотрели ваши заявки вместе. Если нет - не беда. На нашем кэмпе ты познакомишься с множеством прекрасных и удивительных людей", chatId, NilKeyboard)
		case v.companion == "":
			v.companion = text
			sendMessageReq("Есть ли у тебя хронические или иные заболевания, ограничивающие физическую активность или особые пищевые привычки?", chatId, NilKeyboard)
		case v.restricts == "":
			v.restricts = text
			sendMessageReq("Как ты узнала/узнал про RTG? Интересно изучить конверсию", chatId, NilKeyboard)
		case v.conversion == "":
			v.conversion = text
			sendMessageReq("Ты подтверждаешь свое согласие с условиями участия в мероприятии(обработка персональных данных и пр.)? Ознакомиться с условиями следует по ссылке https://bit.ly/2KFSYl8", chatId, YesNoKeyboard)
			// да/нет
		case v.agreement == "":
			v.agreement = text
			sendMessageReq("Контакт родственника/близкого на случай ЧП (имя, как связаться). Только на крайний случай", chatId, NilKeyboard)
		case v.relative == "":
			v.relative = text
			sendMessageReq("Завка отправлена", chatId, NilKeyboard)
			delete(openForms, chatId)

			//todo:проверка формы перед отправкой с возможностью редактирования
			//отправить в гугл форму
		}

	}

	fmt.Println("reply sent")
}

func continiousFetchExpedition() {
	c := time.Tick(5 * time.Minute)
	for now := range c {
		fmt.Println("Expedition list fetched", now)
		expeditionList = parsing.FetchExpeditionsFromUrl(expeditionUrl)
	}

}

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
func setupYesNoKeyboard() *inlineKeyboardMarkup {

	buttonYes := inlineKeyboardButton {
		Text:         "Yes",
		CallbackData: "Yes",
	}
	buttonNo := inlineKeyboardButton {
		Text:         "No",
		CallbackData: "No",
	}
	//test := [][]int{{1,2,3},{1,2,3}}
	AllButtons := inlineKeyboardMarkup{[][]inlineKeyboardButton {{buttonYes}, {buttonNo}}}
	return &AllButtons
}

func sendMessageReq(text string, chatID int64, keyboard int) error {

	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
 	}

	if keyboard == YesNoKeyboard {
		println("hihi\n")
		reqBody.ReplyMarkup = setupYesNoKeyboard()
	}

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	// Send a post request with your token
	var sendMessageUrl = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiString)
	res, err := http.Post(sendMessageUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}
	return nil
}

func find(chatID int64) error {
	var text string
	if expeditionList == nil {
		err := sendMessageReq("Сейчас нет доступных экспедиций", chatID, NilKeyboard)
		if err != nil {
			return err
		}
		return nil
	}
	for _, val := range expeditionList {
		text = "<b>" + val.Name + "</b> \n"
		text += val.Place + "\n"
		text += val.Link + "\n"
		text += "\n"

		err := sendMessageReq(text, chatID, ChooseExpeditionKeyboard)
		if err != nil {
			return err
		}
	}
	return nil
}

func ask(chatID int64, userName string) error {
	fmt.Println("hi" + userName)
	delete(openForms, chatID)
	openForms[chatID] = &Form{tg_name: userName}

	err := sendMessageReq("Куда идем?", chatID, NilKeyboard)
	if err != nil {
		return err
	}

	return nil
}

// func ask_expediton(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, form *Form, expeditionList []parsing.Expedition) {
// 	msg.Text = "Куда идем?"
// 	msg.ReplyMarkup = setupKeyboard(expeditionList)

// 	if _, err := bot.Send(msg); err != nil {
// 		log.Panic(err)
// 	}

// 	for update := range updates {
// 		if update.CallbackQuery != nil {
// 			// Respond to the callback query, telling Telegram to show the user
// 			// a message with the data received.
// 			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
// 			if _, err := bot.Request(callback); err != nil {
// 				panic(err)
// 			}
// 			form.expediton = callback.Text
// 			// And finally, send a message containing the data received.
// 			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
// 			if _, err := bot.Send(msg); err != nil {
// 				panic(err)
// 			}
// 			break
// 		}
// 	}
// }

// func ask(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, form *Form, q string) {
// 	msg.Text = q
// 	if _, err := bot.Send(msg); err != nil {
// 		log.Panic(err)
// 	}

// 	for update := range updates {
// 		if update.Message != nil {
// 			//это можно сделать без повторения списков вопросов? можно. Как?
// 			switch q {
// 			case "Имя и фамилия:":
// 				form.name = update.Message.Text
// 			case "Электронный адрес:":
// 				form.email = update.Message.Text
// 			case "Возраст:":
// 				form.age = update.Message.Text
// 			}
// 			break
// 		}
// 	}
// }

func setupKeyboard(expeditionList []parsing.Expedition) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(expeditionList))

	for _, expedition := range expeditionList {
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(expedition.Name, expedition.Place)))
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	return numericKeyboard
}

// func find(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
// 	expList := parsing.FetchExpeditionsFromUrl("https://russiantravelgeek.com/expeditions/") //убрать
// 	msg.ParseMode = "HTML"

// 	//msg.DisableWebPagePreview = false
// 	for _, val := range expList {
// 		msg.Text = "<b>" + val.Name + "</b> \n"
// 		msg.Text += val.Place + "\n"
// 		msg.Text += val.Link + "\n"
// 		msg.Text += "\n"

// 		if _, err := bot.Send(msg); err != nil {
// 			log.Panic(err)
// 		}
// 	}
// }

// fmt.Println("Started")
// bot, err := tgbotapi.NewBotAPI(apiString)
// if err != nil {
// 	log.Panic(err)
// }
// //bot.Debug = true
// //log.Printf("Authorized on account %s", bot.Self.UserName)

// u := tgbotapi.NewUpdate(0)
// u.Timeout = 60

// expeditionList := parsing.FetchExpeditionsFromUrl(expeditionUrl)
// go continiousParseExpedition(expeditionList)

// var form Form

// updates := bot.GetUpdatesChan(u)
// for update := range updates {
// 	if update.Message == nil { // ignore any non-Message updates
// 		continue
// 	}
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
// 	if !update.Message.IsCommand() { // ignore any non-command Messages
// 		msg.Text = "Чтобы посмотреть список экспедиций, нажмите /find \nЧтобы подать заявку на экспедицию, нажмите /go"
// 		if _, err := bot.Send(msg); err != nil {
// 			log.Panic(err)
// 		}
// 		continue
// 	}

// 	switch update.Message.Command() {
// 	case "find":
// 		find(msg, bot)
// 	case "go":
// 		form.tg_name = update.Message.From.UserName
// 		ask_expediton(msg, bot, updates, &form, expeditionList)
// 		ask(msg, bot, updates, &form, "Имя и фамилия:")
// 		ask(msg, bot, updates, &form, "Электронный адрес:")
// 		ask(msg, bot, updates, &form, "Возраст:")
// 		check_info(msg, bot, form)
// 		//send_info(form)
// 		//https://t.me/+aCYdv4e0hmw0NWUy
// 	default:
// 		msg.Text = "I don't know that command"
// 	}
// }
