package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("8364134131:AAFZCyJWV_SJlbEEc0Uvhwl_wAaudRWFAfQ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot authorized as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		// Пришло обычное сообщение — показываем кнопки
		if update.Message != nil {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("💣 Bomb", "bomb"),
					tgbotapi.NewInlineKeyboardButtonData("🪙 Token", "token"),
				),
			)

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Выбери действие:",
			)
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}

		// Нажатие на кнопку
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			var text string

			switch data {
			case "bomb":
				text = "Нажата кнопка 💣 Bomb"
			case "token":
				text = "Нажата кнопка 🪙 Token"
			default:
				text = "Неизвестная кнопка"
			}

			// Отправляем сообщение
			msg := tgbotapi.NewMessage(chatID, text)
			bot.Send(msg)

			// Обязательно подтверждаем callback
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.Request(callback)
		}
	}
}
