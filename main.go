package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Берём токен из переменной окружения
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot authorized as @%s", bot.Self.UserName)
	log.Println("Bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		// Любое сообщение — показываем кнопки
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

		// Обработка нажатия кнопок
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

			msg := tgbotapi.NewMessage(chatID, text)
			bot.Send(msg)

			// Подтверждаем callback (обязательно)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.Request(callback)
		}
	}
}
