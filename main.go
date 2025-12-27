package main

import (
	"log"
	"os"
	table "tacticus-tg-bot/result-drawer"
	tacticus "tacticus-tg-bot/tacticus-wrapper"
	yandex "tacticus-tg-bot/yandex-tables-wrapper"
	"time"

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

		// Любое сообщение показываем кнопки
		if update.Message != nil {
			getBotMenu(update.Message.Chat.ID, bot)
		}

		// Обработка нажатия кнопок
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			switch data {
			case "bomb":
				runFetchBombProcess(chatID, bot)
			case "token":
			default:
			}

			getBotMenu(chatID, bot)

			// Подтверждаем callback (обязательно)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			_, _ = bot.Request(callback)
		}
	}
}

func runFetchBombProcess(chatID int64, bot *tgbotapi.BotAPI) bool {
	// Канал для остановки спиннера
	done := make(chan bool)

	// Отправляем начальное сообщение
	msg := tgbotapi.NewMessage(chatID, "🔄 Работаю")
	sentMsg, _ := bot.Send(msg)

	go func() {
		spinner := []string{".", "..", "...", "...."}
		i := 0
		for {
			select {
			case <-done:
				deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
				_, _ = bot.Request(deleteMsg)
				return // остановка
			default:
				edit := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, "🔄 Работаю"+spinner[i%len(spinner)])
				_, _ = bot.Send(edit)
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	imagePath, err := fetchPlayersBombsInfo()
	done <- true
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, imagePath))
		return false
	}

	// Отправляем картинку в Telegram
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	_, _ = bot.Send(photoMsg)
	_ = os.Remove(imagePath)

	return true
}

func getBotMenu(chatID int64, bot *tgbotapi.BotAPI) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💣 Bomb", "bomb"),
			tgbotapi.NewInlineKeyboardButtonData("🪙 Token", "token"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Выбери действие:")
	msg.ReplyMarkup = keyboard
	_, _ = bot.Send(msg)
}

func fetchPlayersBombsInfo() (string, error) {

	link, err := yandex.GetYandexDownloadLink("https://disk.yandex.ru/i/mpnEtY1HjtAg8Q")
	if err != nil {
		log.Println("Ошибка получения ссылки:", err)
		return "❌ Ошибка получения ссылки", err
	}

	file, err := yandex.LoadXLSX(link)
	if err != nil {
		log.Println("Ошибка загрузки XLSX:", err)
		return "❌ Ошибка загрузки XLSX", err
	}

	players, err := yandex.ReadPlayersFromXLSX(file)
	if err != nil {
		log.Println("Ошибка чтения игроков:", err)
		return "❌ Ошибка чтения игроков", err
	}

	resultPlayers := tacticus.FetchPlayers(players)
	imagePath := "table.png"
	err = table.DrawImageWithTables(resultPlayers, imagePath)
	if err != nil {
		log.Println("Ошибка генерации картинки:", err)
		return "❌ Ошибка генерации картинки", err
	}

	return imagePath, nil
}
